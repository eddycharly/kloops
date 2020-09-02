package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	kindBitbucketServer = "bitbucketserver"
)

// TokenFunc generates a credentials token
type TokenFunc func() []byte

// Client represents a git client
type Client interface {
	Clean() error
	Clone(repo string, base string, user string, token TokenFunc) (*Repo, error)
}

// client can clone repos. It keeps a local cache, so successive clones of the
// same repo should be quick. Launch with NewClient. Be sure to clean it up.
type client struct {
	// logger will be used to log git operations and must be set.
	logger logr.Logger
	// dir is the location of the git cache.
	dir string
	// git is the path to the git binary.
	git string
	// The mutex protects repoLocks which protect individual repos. This is
	// necessary because Clone calls for the same repo are racy. Rather than
	// one lock for all repos, use a lock per repo.
	// Lock with Client.lockRepo, unlock with Client.unlockRepo.
	rlm       sync.Mutex
	repoLocks map[string]*sync.Mutex
	credLock  sync.RWMutex
}

// Clean removes the local repo cache. The Client is unusable after calling.
func (c *client) Clean() error {
	return os.RemoveAll(c.dir)
}

// NewClient returns a client that talks to git. It will fail if git is not in the PATH.
func NewClient() (Client, error) {
	t, err := ioutil.TempDir("", "git")
	if err != nil {
		return nil, err
	}
	return NewClientWithDir(t)
}

// NewClientWithDir uses an existing directory for creating the client.
func NewClientWithDir(dir string) (Client, error) {
	g, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}
	return &client{
		logger:    ctrl.Log.WithName("git"),
		dir:       dir,
		git:       g,
		repoLocks: make(map[string]*sync.Mutex),
	}, nil
}

func (c *client) lockRepo(repo string) {
	c.rlm.Lock()
	if _, ok := c.repoLocks[repo]; !ok {
		c.repoLocks[repo] = &sync.Mutex{}
	}
	m := c.repoLocks[repo]
	c.rlm.Unlock()
	m.Lock()
}

func (c *client) unlockRepo(repo string) {
	c.rlm.Lock()
	defer c.rlm.Unlock()
	c.repoLocks[repo].Unlock()
}

// Clone clones a repository. Pass the full repository name, such as
// "kubernetes/test-infra" as the repo.
// This function may take a long time if it is the first time cloning the repo.
// In that case, it must do a full git mirror clone. For large repos, this can
// take a while. Once that is done, it will do a git fetch instead of a clone,
// which will usually take at most a few seconds.
func (c *client) Clone(repo string, base string, user string, token TokenFunc) (*Repo, error) {
	c.lockRepo(repo)
	defer c.unlockRepo(repo)

	pass := string(token())
	if user != "" && pass != "" {
		host := gitHost(base)
		base = fmt.Sprintf("https://%s:%s@%s", user, pass, host)
	}
	cache := filepath.Join(c.dir, repo) + ".git"
	if _, err := os.Stat(cache); os.IsNotExist(err) {
		// Cache miss, clone it now.
		c.logger.WithValues("repo", repo).Info("Cloning for the first time")
		if err := os.Mkdir(filepath.Dir(cache), os.ModePerm); err != nil && !os.IsExist(err) {
			return nil, err
		}
		prefix := ""
		repoText := repo
		// if c.gitKind == kindBitbucketServer {
		// 	prefix = "scm/"
		// 	idx := strings.Index(repo, "/")

		// 	// to clone on bitbucket we need to lower case the projectKey owner
		// 	if idx > 0 {
		// 		repoText = fmt.Sprintf("%s/%s", strings.ToLower(repo[0:idx]), repo[idx+1:])
		// 	}
		// }
		remote := fmt.Sprintf("%s/%s%s", base, prefix, repoText)
		if b, err := retryCmd(c.logger, "", c.git, "clone", "--mirror", remote, cache); err != nil {
			return nil, fmt.Errorf("git cache clone error: %v. output: %s", err, string(b))
		}
	} else if err != nil {
		return nil, err
	} else {
		// Cache hit. Do a git fetch to keep updated.
		c.logger.WithValues("repo", repo).Info("Fetching")
		if b, err := retryCmd(c.logger, cache, c.git, "fetch"); err != nil {
			return nil, fmt.Errorf("git fetch error: %v. output: %s", err, string(b))
		}
	}
	t, err := ioutil.TempDir("", "git")
	if err != nil {
		return nil, err
	}
	b, err := exec.Command(c.git, "clone", cache, t).CombinedOutput() // #nosec
	if err != nil {
		return nil, fmt.Errorf("git repo clone error: %v. output: %s", err, string(b))
	}
	return &Repo{
		dir:    t,
		logger: c.logger.WithValues("repo", repo),
		git:    c.git,
		base:   base,
		repo:   repo,
		user:   user,
		pass:   pass,
	}, nil
}

func gitHost(s string) string {
	u, err := url.Parse(s)
	if err == nil {
		return u.Host
	}
	return strings.TrimPrefix(s, "https://")
}

// Repo is a clone of a git repository. Launch with Client.Clone, and don't
// forget to clean it up after.
type Repo struct {
	// dir is the location of the git repo.
	dir string
	// git is the path to the git binary.
	git string
	// base is the base path for remote git fetch calls.
	base string
	// repo is the full repo name: "org/repo".
	repo string
	// user is used for pushing to the remote repo.
	user string
	// pass is used for pushing to the remote repo.
	pass string
	// logger
	logger logr.Logger
}

// Clean deletes the repo. It is unusable after calling.
func (r *Repo) Clean() error {
	return os.RemoveAll(r.dir)
}

func (r *Repo) gitCommand(arg ...string) *exec.Cmd {
	cmd := exec.Command(r.git, arg...) // #nosec
	cmd.Dir = r.dir
	return cmd
}

// Checkout runs git checkout.
func (r *Repo) Checkout(commitlike string) error {
	r.logger.WithValues("commitlike", commitlike).Info("Checkout")
	co := r.gitCommand("checkout", commitlike)
	if b, err := co.CombinedOutput(); err != nil {
		return fmt.Errorf("error checking out %s: %v. output: %s", commitlike, err, string(b))
	}
	return nil
}

// RevParse runs git rev-parse.
func (r *Repo) RevParse(commitlike string) (string, error) {
	r.logger.WithValues("commitlike", commitlike).Info("RevParse")
	b, err := r.gitCommand("rev-parse", commitlike).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error rev-parsing %s: %v. output: %s", commitlike, err, string(b))
	}
	return string(b), nil
}

// CheckoutNewBranch creates a new branch and checks it out.
func (r *Repo) CheckoutNewBranch(branch string) error {
	r.logger.WithValues("branch", branch).Info("Launch and checkout")
	co := r.gitCommand("checkout", "-b", branch)
	if b, err := co.CombinedOutput(); err != nil {
		return fmt.Errorf("error checking out %s: %v. output: %s", branch, err, string(b))
	}
	return nil
}

// Merge attempts to merge commitlike into the current branch. It returns true
// if the merge completes. It returns an error if the abort fails.
func (r *Repo) Merge(commitlike string) (bool, error) {
	r.logger.WithValues("commitlike", commitlike).Info("Merging")
	co := r.gitCommand("merge", "--no-ff", "--no-stat", "-m merge", commitlike)

	b, err := co.CombinedOutput()
	if err == nil {
		return true, nil
	}
	r.logger.WithValues("commitlike", commitlike, "output", string(b)).Error(err, "Merge failed")

	if b, err := r.gitCommand("merge", "--abort").CombinedOutput(); err != nil {
		return false, fmt.Errorf("error aborting merge for commitlike %s: %v. output: %s", commitlike, err, string(b))
	}

	return false, nil
}

// Am tries to apply the patch in the given path into the current branch
// by performing a three-way merge (similar to git cherry-pick). It returns
// an error if the patch cannot be applied.
func (r *Repo) Am(path string) error {
	r.logger.WithValues("path", path).Info("Applying")
	co := r.gitCommand("am", "--3way", path)
	b, err := co.CombinedOutput()
	if err == nil {
		return nil
	}
	output := string(b)
	r.logger.WithValues("output", output).Error(err, "Patch apply failed")
	if b, abortErr := r.gitCommand("am", "--abort").CombinedOutput(); abortErr != nil {
		r.logger.WithValues("output", string(b)).Info("Aborting patch apply failed")
	}
	applyMsg := "The copy of the patch that failed is found in: .git/rebase-apply/patch"
	if strings.Contains(output, applyMsg) {
		i := strings.Index(output, applyMsg)
		err = fmt.Errorf("%s", output[:i])
	}
	return err
}

// Push pushes over https to the provided owner/repo#branch using a password
// for basic auth.
func (r *Repo) Push(repo, branch string) error {
	if r.user == "" || r.pass == "" {
		return errors.New("cannot push without credentials - configure your git client")
	}
	r.logger.WithValues("user", r.user, "repo", repo, "branch", branch).Info("Pushing")
	host := gitHost(r.base)
	remote := fmt.Sprintf("https://%s:%s@%s/%s/%s", r.user, r.pass, host, r.user, repo)
	co := r.gitCommand("push", remote, branch)
	_, err := co.CombinedOutput()
	return err
}

// CheckoutPullRequest does exactly that.
func (r *Repo) CheckoutPullRequest(number int) error {
	r.logger.WithValues("repo", r.repo, "number", number).Info("Fetching and checking out")
	if b, err := retryCmd(r.logger, r.dir, r.git, "fetch", r.base+"/"+r.repo, fmt.Sprintf("pull/%d/head:pull%d", number, number)); err != nil {
		return fmt.Errorf("git fetch failed for PR %d: %v. output: %s", number, err, string(b))
	}
	co := r.gitCommand("checkout", fmt.Sprintf("pull%d", number))
	if b, err := co.CombinedOutput(); err != nil {
		return fmt.Errorf("git checkout failed for PR %d: %v. output: %s", number, err, string(b))
	}
	return nil
}

// Config runs git config.
func (r *Repo) Config(key, value string) error {
	r.logger.WithValues("key", key, "value", value).Info("Running git config")
	if b, err := r.gitCommand("config", key, value).CombinedOutput(); err != nil {
		return fmt.Errorf("git config %s %s failed: %v. output: %s", key, value, err, string(b))
	}
	return nil
}

// retryCmd will retry the command a few times with backoff. Use this for any commands that will be talking to git, such as clones or fetches.
func retryCmd(logger logr.Logger, dir, cmd string, arg ...string) ([]byte, error) {
	var b []byte
	var err error
	sleepyTime := time.Second
	for i := 0; i < 3; i++ {
		c := exec.Command(cmd, arg...) // #nosec
		c.Dir = dir
		b, err = c.CombinedOutput()
		if err != nil {
			logger.WithValues("cmd", cmd, "arg", arg, "output", string(b)).Error(err, "Running command returned error")
			time.Sleep(sleepyTime)
			sleepyTime *= 2
			continue
		}
		break
	}
	return b, err
}
