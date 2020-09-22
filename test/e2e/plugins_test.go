package e2e

import (
	"context"
	"testing"

	"github.com/eddycharly/kloops/pkg/git"
	"github.com/jenkins-x/go-scm/scm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	return Describe("Cat plugin", func() {
		var (
			bot        client
			repo       *scm.Repository
			localClone *git.Repo
		)
		AfterEach(func() {
			bot.scmClient.Repositories.Delete(context.Background(), repo.FullName)
		})
		BeforeEach(func() {
			hmacToken := SetupHmac()
			gitServerURL := GetGitServerURL()
			bot = SetupBotClient()
			gitClient := SetupGitClient(gitServerURL)
			repo, localClone = SetupRepo(gitServerURL, bot.scmClient, gitClient)
			SetupConfig(repo, ns, hmacToken, GetPrimarySCMToken)
			SetupWebhook(bot.scmClient, repo, hmacToken)
		})
		Describe("With pull requests", func() {
			var (
				pr *scm.PullRequest
			)
			BeforeEach(func() {
				pr = SetupPullRequest(bot.scmClient, repo, localClone, prBranch)
			})
			It("Checks /meow works", func() {
				var (
					err error
				)
				By("adding '/meow' comment", func() {
					err = bot.providerClient.PullRequests.CreateComment(repo.FullName, pr.Number, "/meow")
					Expect(err).ShouldNot(HaveOccurred())
					err = WaitForPullRequestComment(bot.providerClient, pr, `!\[cat image\]`)
					Expect(err).NotTo(HaveOccurred())
				})
			})
			It("Checks /meowvie works", func() {
				var (
					err error
				)
				By("adding '/meowvie' comment", func() {
					err = bot.providerClient.PullRequests.CreateComment(repo.FullName, pr.Number, "/meowvie")
					Expect(err).ShouldNot(HaveOccurred())
					err = WaitForPullRequestComment(bot.providerClient, pr, `!\[cat image\]`)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
		Describe("With issues", func() {
			var (
				issue *scm.Issue
			)
			BeforeEach(func() {
				issue = SetupIssue(bot.scmClient, repo)
			})
			It("Checks /meow works", func() {
				var (
					err error
				)
				By("adding '/meow' comment", func() {
					err = bot.providerClient.PullRequests.CreateComment(repo.FullName, issue.Number, "/meow")
					Expect(err).ShouldNot(HaveOccurred())
					err = WaitForIssueComment(bot.providerClient, repo, issue, `!\[cat image\]`)
					Expect(err).NotTo(HaveOccurred())
				})
			})
			It("Checks /meowvie works", func() {
				var (
					err error
				)
				By("adding '/meowvie' comment", func() {
					err = bot.providerClient.PullRequests.CreateComment(repo.FullName, issue.Number, "/meowvie")
					Expect(err).ShouldNot(HaveOccurred())
					err = WaitForIssueComment(bot.providerClient, repo, issue, `!\[cat image\]`)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})
}
