package e2e

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/jenkins-x/go-scm/scm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ns             = "tools"
	prBranch       = "for-pr"
	defaultContext = "pr-build"
)

func TestTekton(t *testing.T) {
	RunWithReporters(t, "KLoops integration")
}

var _ = CatTests()

func CatTests() bool {
	return Describe("Cat plugin support", func() {
		var (
			hmacToken      string
			gitClient      git.Client
			scmClient      *scm.Client
			spc            *scmprovider.Client
			approverClient *scm.Client
			gitServerURL   string
			repo           *scm.Repository
			repoFullName   string
			localClone     *git.Repo
		)
		// AfterEach(func() {
		// 	gitClient.Clean()
		// })
		BeforeEach(func() {
			var err error
			By("creating HMAC token")
			hmacToken, err = CreateHMACToken()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(hmacToken).ShouldNot(BeEmpty())

			By("creating primary SCM client")
			scmClient, spc, gitServerURL, err = CreateSCMClient(GetBotName, GetPrimarySCMToken)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(scmClient).ShouldNot(BeNil())
			// Expect(spc).ShouldNot(BeNil())
			Expect(gitServerURL).ShouldNot(BeEmpty())

			By("creating approver SCM client")
			approverClient, _, _, err = CreateSCMClient(GetApproverName, GetApproverSCMToken)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(approverClient).ShouldNot(BeNil())
			// Expect(approverSpc).ShouldNot(BeNil())

			By("creating git client")
			gitClient, err = CreateGitClient(gitServerURL, GetBotName, GetPrimarySCMToken)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(gitClient).ShouldNot(BeNil())

			By("creating repository")
			repo, localClone, err = CreateBaseRepository(GetBotName(), GetApproverName(), gitServerURL, "gitea", GetPrimarySCMToken, scmClient, gitClient)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(repo).ShouldNot(BeNil())
			Expect(localClone).ShouldNot(BeNil())
			repoFullName = fmt.Sprintf("%s/%s", repo.Namespace, repo.Name)
			fmt.Println(repoFullName)
			By(fmt.Sprintf("adding %s to new repository", GetApproverName()))
			err = AddCollaborator(GetApproverName(), repo, scmClient, approverClient)
			Expect(err).ShouldNot(HaveOccurred())

			token, err := GetPrimarySCMToken()
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

			By(fmt.Sprintf("setting up webhooks for %s", repo.Clone))
			err = CreateWebHook(scmClient, repo, hmacToken)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Checks /meow works", func() {
			var (
				err error
				pr  *scm.PullRequest
			)
			By("cloning, creating the new branch, and pushing it", func() {
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
			})
			By("creating a pull request", func() {
				prInput := &scm.PullRequestInput{
					Title: "KLoops Test PR",
					Head:  prBranch,
					Base:  "master",
					Body:  "Test PR for KLoops",
				}
				pr, _, err = scmClient.PullRequests.Create(context.Background(), repoFullName, prInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(pr).ShouldNot(BeNil())
			})
			By("adding '/meow' comment", func() {
				err = spc.PullRequests.CreateComment(repo.FullName, pr.Number, "/meow")
				Expect(err).ShouldNot(HaveOccurred())
				err = WaitForPullRequestComment(spc, pr, `!\[cat image\]`)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
}
