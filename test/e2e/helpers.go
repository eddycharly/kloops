package e2e

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/eddycharly/kloops/pkg/repoowners"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/onsi/ginkgo"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	gr "github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

const (
	primarySCMTokenEnvVar  = "E2E_PRIMARY_SCM_TOKEN" /* #nosec */
	primarySCMUserEnvVar   = "E2E_PRIMARY_SCM_USER"
	approverSCMTokenEnvVar = "E2E_APPROVER_SCM_TOKEN" /* #nosec */
	approverSCMUserEnvVar  = "E2E_APPROVER_SCM_USER"
	hmacTokenEnvVar        = "E2E_HMAC_TOKEN" /* #nosec */
	gitServerEnvVar        = "E2E_GIT_SERVER"
	gitKindEnvVar          = "E2E_GIT_KIND"
	kloopsURLEnvVar        = "E2E_KLOOPS_URL"
	baseRepoName           = "lh-e2e-test"
)

// RunWithReporters runs a suite with better logging and gathering of test results
func RunWithReporters(t *testing.T, suiteID string) {
	reportsDir := os.Getenv("REPORTS_DIR")
	if reportsDir == "" {
		reportsDir = filepath.Join("../", "build", "reports")
	}
	err := os.MkdirAll(reportsDir, 0700)
	if err != nil {
		t.Errorf("cannot create %s because %v", reportsDir, err)
	}
	reporters := make([]ginkgo.Reporter, 0)

	slowSpecThresholdStr := os.Getenv("SLOW_SPEC_THRESHOLD")
	if slowSpecThresholdStr == "" {
		slowSpecThresholdStr = "50000"
		_ = os.Setenv("SLOW_SPEC_THRESHOLD", slowSpecThresholdStr)
	}
	slowSpecThreshold, err := strconv.ParseFloat(slowSpecThresholdStr, 64)
	if err != nil {
		panic(err.Error())
	}
	ginkgoconfig.DefaultReporterConfig.SlowSpecThreshold = slowSpecThreshold
	ginkgoconfig.DefaultReporterConfig.Verbose = testing.Verbose()
	reporters = append(reporters, gr.NewJUnitReporter(filepath.Join(reportsDir, fmt.Sprintf("%s.junit.xml", suiteID))))
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, fmt.Sprintf("KLoops E2E tests: %s", suiteID), reporters)
}

// CreateHMACToken creates an HMAC token for use in webhooks, defaulting to the E2E_HMAC_TOKEN env var if set
func CreateHMACToken() (string, error) {
	fromEnv := os.Getenv(hmacTokenEnvVar)
	if fromEnv != "" {
		return fromEnv, nil
	}
	src := rand.New(rand.NewSource(time.Now().UnixNano())) /* #nosec */
	b := make([]byte, 21)                                  // can be simplified to n/2 if n is always even
	if _, err := src.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b)[:41], nil
}

func GetGitServerURL() string {
	return os.Getenv(gitServerEnvVar)
}

// CreateSCMClient takes functions that return the username and token to use, and creates the scm.Client and Lighthouse SCM client
func CreateSCMClient(userFunc func() string, tokenFunc func() (string, error)) (*scm.Client, *scmprovider.Client, string, error) {
	kind := GitKind()
	serverURL := GetGitServerURL()

	token, err := tokenFunc()
	if err != nil {
		return nil, nil, "", err
	}

	client, err := factory.NewClient(kind, serverURL, token)

	// util.AddAuthToSCMClient(client, token, false)

	spc := scmprovider.NewClient(client)

	return client, &spc, serverURL, err
}

// GetBotName gets the bot user name
func GetBotName() string {
	botName := os.Getenv(primarySCMUserEnvVar)
	if botName == "" {
		botName = "kloops-bot"
	}
	return botName
}

// GetPrimarySCMToken gets the token used by the bot/primary user
func GetPrimarySCMToken() (string, error) {
	return getSCMToken(primarySCMTokenEnvVar, GitKind())
}

// GetApproverName gets the approver user's username
func GetApproverName() string {
	botName := os.Getenv(approverSCMUserEnvVar)
	if botName == "" {
		botName = "kloops-bot"
	}
	return botName
}

// GetApproverSCMToken gets the token used by the approver
func GetApproverSCMToken() (string, error) {
	return getSCMToken(approverSCMTokenEnvVar, GitKind())
}

// GetWebhookURL the webhooks url
func GetWebhookURL(repo *scm.Repository) string {
	return fmt.Sprintf("%s/%s/%s", os.Getenv(kloopsURLEnvVar), "tools", repo.Name)
}

func getSCMToken(envName, gitKind string) (string, error) {
	value := os.Getenv(envName)
	if value == "" {
		return value, fmt.Errorf("No token available for git kind %s at environment variable $%s", gitKind, envName)
	}
	return value, nil
}

// GitKind returns the git provider flavor being used
func GitKind() string {
	kind := os.Getenv(gitKindEnvVar)
	if kind == "" {
		kind = "gitea"
	}
	return kind
}

// CreateGitClient creates the git client used for cloning and making changes to the test repository
func CreateGitClient(gitServerURL string, userFunc func() string, tokenFunc func() (string, error)) (git.Client, error) {
	gitClient, err := git.NewClientWithDir("/Users/charlesbreteche/Documents/dev/eddycharly/kloops/tempo")
	// gitClient, err := git.NewClient()
	if err != nil {
		return nil, err
	}
	// token, err := tokenFunc()
	// if err != nil {
	// 	return nil, err
	// }
	// gitClient.SetCredentials(userFunc(), func() []byte {
	// 	return []byte(token)
	// })

	return gitClient, nil
}

// CreateBaseRepository creates the repository that will be used for tests
func CreateBaseRepository(botUser, approver, gitServerURL, user string, tokenFunc func() (string, error), botClient *scm.Client, gitClient git.Client) (*scm.Repository, *git.Repo, error) {
	repoName := baseRepoName + "-" + strconv.FormatInt(ginkgo.GinkgoRandomSeed(), 10)

	input := &scm.RepositoryInput{
		Name:    repoName,
		Private: true,
	}

	repo, _, err := botClient.Repositories.Create(context.Background(), input)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create repository")
	}

	// Sleep 5 seconds to ensure repository exists enough to be pushed to.
	time.Sleep(5 * time.Second)
	token, err := tokenFunc()
	if err != nil {
		return nil, nil, err
	}

	r, err := gitClient.Clone(repo.FullName, gitServerURL, user, func() []byte { return []byte(token) })
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not clone %s", repo.FullName)
	}
	err = r.CheckoutNewBranch("master")
	if err != nil {
		return nil, nil, err
	}

	baseScriptFile := filepath.Join("test_data", "baseRepoScript.sh")
	baseScript, err := ioutil.ReadFile(baseScriptFile) /* #nosec */

	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to read %s", baseScriptFile)
	}
	fmt.Println(r.Dir())
	scriptOutputFile := filepath.Join(r.Dir(), "script.sh")
	err = ioutil.WriteFile(scriptOutputFile, baseScript, 0600)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't write to %s", scriptOutputFile)
	}

	ExpectCommandExecution(r.Dir(), 1, 0, "git", "add", scriptOutputFile)

	owners := repoowners.SimpleConfig{
		Config: repoowners.Config{
			Approvers: []string{botUser, approver},
			Reviewers: []string{botUser, approver},
		},
	}

	ownersFile := filepath.Join(r.Dir(), "OWNERS")
	ownersYaml, err := yaml.Marshal(owners)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't marshal OWNERS yaml")
	}

	err = ioutil.WriteFile(ownersFile, ownersYaml, 0600)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't write to %s", ownersFile)
	}
	ExpectCommandExecution(r.Dir(), 1, 0, "git", "add", ownersFile)
	ExpectCommandExecution(r.Dir(), 1, 0, "git", "commit", "-a", "-m", "Initial commit of functioning script and OWNERS")

	err = r.Push(repo.Name, "master")
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to push to %s", repo.Clone)
	}

	return repo, r, nil
}

// ExpectCommandExecution performs the given command in the current work directory and asserts that it completes successfully
func ExpectCommandExecution(dir string, commandTimeout time.Duration, exitCode int, c string, args ...string) {
	f := func() error {
		command := exec.Command(c, args...) /* #nosec */
		command.Dir = dir
		session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
		session.Wait(10 * time.Second * commandTimeout)
		gomega.Eventually(session).Should(gexec.Exit(exitCode))
		return err
	}
	err := retryExponentialBackoff(1, f)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
}

// retryExponentialBackoff retries the given function up to the maximum duration
func retryExponentialBackoff(maxDuration time.Duration, f func() error) error {
	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.MaxElapsedTime = maxDuration
	exponentialBackOff.Reset()
	err := backoff.Retry(f, exponentialBackOff)
	return err
}

// AddCollaborator adds the approver user to the repo
func AddCollaborator(approver string, repo *scm.Repository, botClient *scm.Client, approverClient *scm.Client) error {
	_, alreadyCollaborator, _, err := botClient.Repositories.AddCollaborator(context.Background(), fmt.Sprintf("%s/%s", repo.Namespace, repo.Name), approver, "admin")
	if alreadyCollaborator {
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "adding %s as collaborator for repo %s/%s", approver, repo.Namespace, repo.Name)
	}

	// Don't bother checking for invites with BitBucket Server
	if GitKind() == "stash" {
		return nil
	}

	// Sleep for a bit
	time.Sleep(15 * time.Second)

	invites, _, err := approverClient.Users.ListInvitations(context.Background())
	if err == scm.ErrNotSupported {
		// Ignore any cases of not supported
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "listing invitations for user %s", approver)
	}
	for _, i := range invites {
		_, err = approverClient.Users.AcceptInvitation(context.Background(), i.ID)
		if err == scm.ErrNotSupported {
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "accepting invitation %d for user %s", i.ID, approver)
		}
	}
	return nil
}

// CreateWebHook creates a webhook on the SCM provider for the repository
func CreateWebHook(scmClient *scm.Client, repo *scm.Repository, hmacToken string) error {
	input := &scm.HookInput{
		Name:         "kloops-test-hook",
		Target:       GetWebhookURL(repo),
		Secret:       hmacToken,
		NativeEvents: []string{"*"},
	}
	if scmClient.Driver.String() == "gitea" {
		input.Events.Issue = true
		input.Events.PullRequest = true
		input.Events.Branch = true
		input.Events.IssueComment = true
		input.Events.PullRequestComment = true
		input.Events.Push = true
		input.Events.ReviewComment = true
		input.Events.Tag = true
	}
	_, _, err := scmClient.Repositories.CreateHook(context.Background(), repo.FullName, input)
	return err
}

// ApplyConfigAndPluginsConfigMaps takes the config and plugins and creates/applies the config maps in the cluster using kubectl
func ApplyConfigAndPluginsConfigMaps(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfig) error {
	// rc, err := json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil).
	//  .Marshal(repoConfig)

	// if err != nil {
	// 	return errors.Wrapf(err, "writing config to YAML")
	// }
	// pc, err := yaml.Marshal(pluginConfig)
	// if err != nil {
	// 	return errors.Wrapf(err, "writing plugins to YAML")
	// }
	tmpDir, err := ioutil.TempDir("", "kubectl")
	if err != nil {
		return errors.Wrapf(err, "creating temp directory")
	}
	defer os.RemoveAll(tmpDir)

	rcFile, err := os.Create(filepath.Join(tmpDir, "rc.yaml"))
	defer rcFile.Close()
	if err != nil {
		return errors.Wrapf(err, "writing config to YAML")
	}
	pcFile, err := os.Create(filepath.Join(tmpDir, "pc.yaml"))
	defer pcFile.Close()
	if err != nil {
		return errors.Wrapf(err, "writing config to YAML")
	}

	err = json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil).Encode(repoConfig, rcFile)
	if err != nil {
		return errors.Wrapf(err, "writing config map to rc.yaml")
	}
	err = json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil).Encode(pluginConfig, pcFile)
	if err != nil {
		return errors.Wrapf(err, "writing plugins map to pc.ayml")
	}

	ExpectCommandExecution(tmpDir, 1, 0, "kubectl", "apply", "-f", "rc.yaml")
	ExpectCommandExecution(tmpDir, 1, 0, "kubectl", "apply", "-f", "pc.yaml")

	return nil
}
