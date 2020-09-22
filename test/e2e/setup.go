package e2e

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/jenkins-x/go-scm/scm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetupHmac() string {
	By("creating HMAC token")
	hmacToken, err := CreateHMACToken()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(hmacToken).ShouldNot(BeEmpty())
	return hmacToken
}

type client struct {
	scmClient      *scm.Client
	providerClient *scmprovider.Client
}

func SetupClient(name string, userFunc func() string, tokenFunc func() (string, error)) client {
	By(fmt.Sprintf("creating %s client", name))
	scmClient, providerClient, gitServerURL, err := CreateSCMClient(userFunc, tokenFunc)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(scmClient).ShouldNot(BeNil())
	Expect(providerClient).ShouldNot(BeNil())
	Expect(gitServerURL).ShouldNot(BeEmpty())
	return client{
		scmClient:      scmClient,
		providerClient: providerClient,
	}
}

func SetupBotClient() client {
	return SetupClient("bot", GetBotName, GetPrimarySCMToken)
}

func SetupApproverClient() client {
	return SetupClient("approver", GetApproverName, GetApproverSCMToken)
}
func SetupGitClient(gitServerURL string) git.Client {
	By("creating git client")
	gitClient, err := CreateGitClient(gitServerURL, GetBotName, GetPrimarySCMToken)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(gitClient).ShouldNot(BeNil())

	return gitClient
}

func SetupRepo(gitServerURL string, scmClient *scm.Client, gitClient git.Client) (*scm.Repository, *git.Repo) {
	By("creating repository")
	repo, localClone, err := CreateBaseRepository(GetBotName(), GetApproverName(), gitServerURL, "gitea", GetPrimarySCMToken, scmClient, gitClient)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(repo).ShouldNot(BeNil())
	Expect(localClone).ShouldNot(BeNil())
	return repo, localClone
}

func SetupConfig(repo *scm.Repository, ns, hmacToken string, tokenFunc func() (string, error)) {
	By("creating config")
	token, err := tokenFunc()
	Expect(err).ShouldNot(HaveOccurred())
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
	err = ApplyConfigAndPluginsConfigMaps(&rc, &pc)
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
