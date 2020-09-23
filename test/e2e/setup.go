package e2e

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/eddycharly/kloops/pkg/repoowners"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetupHmac() string {
	hmacToken, err := CreateHMACToken()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(hmacToken).ShouldNot(BeEmpty())
	return hmacToken
}

func SetupConfig(repo *scm.Repository, ns, hmacToken string, token string) {
	rc := v1alpha1.RepoConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       "RepoConfig",
			APIVersion: "config.kloops.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      repo.Name,
			Namespace: ns,
		},
		Spec: v1alpha1.RepoConfigSpec{
			Gitea: &v1alpha1.GiteaRepo{
				Owner:     repo.Namespace,
				Repo:      repo.Name,
				ServerURL: "http://gitea-http.tools.svc.cluster.local:3000",
				HmacToken: v1alpha1.Secret{
					Value: hmacToken,
				},
				Token: v1alpha1.Secret{
					Value: token,
				},
			},
			BotName: GetBotName(),
			PluginConfig: v1alpha1.RepoPluginConfig{
				Ref: "test",
				Plugins: []string{
					"cat",
				},
			},
		},
	}
	pc := v1alpha1.PluginConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       "PluginConfig",
			APIVersion: "config.kloops.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "test",
			Namespace: ns,
		},
		Spec: v1alpha1.PluginConfigSpec{
			Label: v1alpha1.Label{
				AdditionalLabels: []string{},
			},
		},
	}
	err := ApplyConfigAndPluginsConfigMaps(&rc, &pc)
	Expect(err).ShouldNot(HaveOccurred())
}

func SetupWebhook(scmClient *scm.Client, repo *scm.Repository, hmacToken string) {
	By(fmt.Sprintf("setting up webhooks for %s", repo.Clone))
	err := CreateWebHook(scmClient, repo, hmacToken)
	Expect(err).ShouldNot(HaveOccurred())
}

func SetupPullRequest(scmClient *scm.Client, repo *scm.Repository, localClone *git.Repo, prBranch string) *scm.PullRequest {
	var err error
	By("cloning, creating the new branch, and pushing it")
	err = localClone.CheckoutNewBranch(prBranch)
	Expect(err).ShouldNot(HaveOccurred())

	newFile := filepath.Join(localClone.Dir(), "README")
	err = ioutil.WriteFile(newFile, []byte("Hello world"), 0600)
	ExpectCommandExecution(localClone.Dir(), 1, 0, "git", "add", newFile)

	changedScriptFile := filepath.Join("test_data", "passingRepoScript.sh")
	changedScript, err := ioutil.ReadFile(changedScriptFile) /* #nosec */
	Expect(err).ShouldNot(HaveOccurred())

	scriptOutputFile := filepath.Join(localClone.Dir(), "script.sh")
	err = ioutil.WriteFile(scriptOutputFile, changedScript, 0600)
	Expect(err).ShouldNot(HaveOccurred())

	ExpectCommandExecution(localClone.Dir(), 1, 0, "git", "commit", "-a", "-m", "Adding for test PR")

	err = localClone.Push(repo.Name, prBranch)
	Expect(err).ShouldNot(HaveOccurred())

	By("creating a pull request")
	prInput := &scm.PullRequestInput{
		Title: "KLoops Test PR",
		Head:  prBranch,
		Base:  "master",
		Body:  "Test PR for KLoops",
	}
	pr, _, err := scmClient.PullRequests.Create(context.Background(), repo.FullName, prInput)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(pr).ShouldNot(BeNil())

	return pr
}

func SetupIssue(scmClient *scm.Client, repo *scm.Repository) *scm.Issue {
	var err error
	By("creating an issue")
	issueInput := &scm.IssueInput{
		Title: "KLoops Test ISSUE",
		Body:  "Test ISSUE for KLoops",
	}
	issue, _, err := scmClient.Issues.Create(context.Background(), repo.FullName, issueInput)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(issue).ShouldNot(BeNil())

	return issue
}

type user struct {
	name      string
	password  string
	client    *scm.Client
	gitClient git.Client
}

func (u user) CreateOrganization(name string) *scm.Organization {
	org, _, err := u.client.Organizations.Create(context.Background(), &scm.OrganizationInput{
		Name: name,
	})
	Expect(err).ShouldNot(HaveOccurred())
	Expect(org).ShouldNot(BeNil())
	return org
}

func (u user) CreateToken(name string) *scm.UserToken {
	token, _, err := u.client.Users.CreateToken(context.Background(), u.name, name)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(token).ShouldNot(BeNil())
	return token
}

func (u user) DeleteToken(id int64) {
	_, err := u.client.Users.DeleteToken(context.Background(), id)
	Expect(err).ShouldNot(HaveOccurred())
}

func (u user) DeleteOrganization(name string) {
	_, err := u.client.Organizations.Delete(context.Background(), name)
	Expect(err).ShouldNot(HaveOccurred())
}

func (u user) CreateRepository(org, name string) *scm.Repository {
	repo, _, err := u.client.Repositories.Create(context.Background(), &scm.RepositoryInput{
		Namespace: org,
		Name:      name,
	})
	Expect(err).ShouldNot(HaveOccurred())
	Expect(repo).ShouldNot(BeNil())
	return repo
}

func (u user) CloneRepository(server, name string) *git.Repo {
	repo, err := u.gitClient.Clone(name, server, u.name, func() []byte { return []byte(u.password) })
	Expect(err).ShouldNot(HaveOccurred())
	Expect(repo).ShouldNot(BeNil())
	return repo
}

func (u user) DeleteRepository(repo string) {
	_, err := u.client.Repositories.Delete(context.Background(), repo)
	Expect(err).ShouldNot(HaveOccurred())
}

func (u user) CreatePullRequest(repo, from, to, title, body string) *scm.PullRequest {
	input := &scm.PullRequestInput{
		Title: title,
		Head:  from,
		Base:  to,
		Body:  body,
	}
	pr, _, err := u.client.PullRequests.Create(context.Background(), repo, input)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(pr).ShouldNot(BeNil())
	return pr
}

func (u user) CreateIssue(repo, title, body string) *scm.Issue {
	input := &scm.IssueInput{
		Title: title,
		Body:  body,
	}
	issue, _, err := u.client.Issues.Create(context.Background(), repo, input)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(issue).ShouldNot(BeNil())
	return issue
}

func (u user) CreateWebhook(repo *scm.Repository, hmacToken string) {
	input := &scm.HookInput{
		Name:   "kloops",
		Target: GetWebhookURL(repo),
		Secret: hmacToken,
		Events: scm.HookEvents{
			Issue:              true,
			PullRequest:        true,
			Branch:             true,
			IssueComment:       true,
			PullRequestComment: true,
			Push:               true,
			ReviewComment:      true,
			Tag:                true,
		},
		NativeEvents: []string{"*"},
	}
	_, _, err := u.client.Repositories.CreateHook(context.Background(), repo.FullName, input)
	Expect(err).ShouldNot(HaveOccurred())
}

func (u user) CreateComment(repo, body string, number int, pr bool) *scm.Comment {
	input := &scm.CommentInput{
		Body: body,
	}
	if pr {
		comment, _, err := u.client.PullRequests.CreateComment(context.Background(), repo, number, input)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(comment).ShouldNot(BeNil())
		return comment
	}
	comment, _, err := u.client.Issues.CreateComment(context.Background(), repo, number, input)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(comment).ShouldNot(BeNil())
	return comment
}

func SetupBasicAuthUser(server, name, password string) user {
	kind := GitKind()
	u, err := url.Parse(server)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(u).ShouldNot(BeNil())
	client, err := factory.NewClientWithBasicAuth(kind, server, name, password)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(client).ShouldNot(BeNil())
	gitClient, err := git.NewClientWithDir("/Users/charlesbreteche/Documents/dev/eddycharly/kloops/tempo")
	Expect(err).ShouldNot(HaveOccurred())
	Expect(gitClient).ShouldNot(BeNil())
	return user{
		name:      name,
		password:  password,
		client:    client,
		gitClient: gitClient,
	}
}

func InitRepository(repo *git.Repo, branch string) {
	err := repo.CheckoutNewBranch(branch)
	Expect(err).ShouldNot(HaveOccurred())

	baseScriptFile := filepath.Join("test_data", "baseRepoScript.sh")
	baseScript, err := ioutil.ReadFile(baseScriptFile) /* #nosec */
	Expect(err).ShouldNot(HaveOccurred())

	scriptOutputFile := filepath.Join(repo.Dir(), "script.sh")
	err = ioutil.WriteFile(scriptOutputFile, baseScript, 0600)
	Expect(err).ShouldNot(HaveOccurred())

	ExpectCommandExecution(repo.Dir(), 1, 0, "git", "add", scriptOutputFile)

	owners := repoowners.SimpleConfig{
		Config: repoowners.Config{
			Approvers: []string{"gitea"},
			Reviewers: []string{"gitea"},
		},
	}

	ownersFile := filepath.Join(repo.Dir(), "OWNERS")
	ownersYaml, err := yaml.Marshal(owners)
	Expect(err).ShouldNot(HaveOccurred())

	err = ioutil.WriteFile(ownersFile, ownersYaml, 0600)
	Expect(err).ShouldNot(HaveOccurred())

	ExpectCommandExecution(repo.Dir(), 1, 0, "git", "add", ownersFile)
	ExpectCommandExecution(repo.Dir(), 1, 0, "git", "commit", "-a", "-m", "Initial commit of functioning script and OWNERS")

	err = repo.Push(repo.Name(), branch)
	Expect(err).ShouldNot(HaveOccurred())
}

func SetupBranch(repo *git.Repo, from, branch string) {
	var err error
	err = repo.Checkout(from)
	Expect(err).ShouldNot(HaveOccurred())
	err = repo.CheckoutNewBranch(branch)
	Expect(err).ShouldNot(HaveOccurred())
	fmt.Println(repo)
	newFile := filepath.Join(repo.Dir(), "README")
	err = ioutil.WriteFile(newFile, []byte("Hello world"), 0600)
	ExpectCommandExecution(repo.Dir(), 1, 0, "git", "add", newFile)

	changedScriptFile := filepath.Join("test_data", "passingRepoScript.sh")
	changedScript, err := ioutil.ReadFile(changedScriptFile) /* #nosec */
	Expect(err).ShouldNot(HaveOccurred())

	scriptOutputFile := filepath.Join(repo.Dir(), "script.sh")
	err = ioutil.WriteFile(scriptOutputFile, changedScript, 0600)
	Expect(err).ShouldNot(HaveOccurred())

	ExpectCommandExecution(repo.Dir(), 1, 0, "git", "commit", "-a", "-m", "Adding for test PR")

	err = repo.Push(repo.Name(), branch)
	Expect(err).ShouldNot(HaveOccurred())
}

func WaitForComment(client *scm.Client, repo string, number int, pr bool, regex string) {
	providerClient := scmprovider.NewClient(client)
	checkComments := func() error {
		var (
			comments []*scm.Comment
			err      error
		)
		if pr {
			comments, err = providerClient.PullRequests.GetComments(repo, number)
		} else {
			comments, err = providerClient.PullRequests.GetComments(repo, number)
		}
		if err != nil {
			return err
		}
		r := regexp.MustCompile(regex)
		for _, comment := range comments {
			if r.MatchString(comment.Body) {
				return nil
			}
		}
		return fmt.Errorf("Failed to find a commeent matching: %s", regex)
	}
	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.MaxElapsedTime = 1 * time.Minute
	exponentialBackOff.MaxInterval = 5 * time.Second
	exponentialBackOff.Reset()
	err := backoff.Retry(checkComments, exponentialBackOff)
	Expect(err).ShouldNot(HaveOccurred())
}
