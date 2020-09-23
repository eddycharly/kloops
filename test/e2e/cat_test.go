package e2e

import (
	"fmt"
	"testing"

	"github.com/eddycharly/kloops/pkg/git"
	"github.com/jenkins-x/go-scm/scm"
	. "github.com/onsi/ginkgo"
)

const (
	ns         = "tools"
	prTitle    = "KLoops Test PR"
	prBody     = "Test PR for KLoops"
	issueTitle = "KLoops Test ISSUE"
	issueBody  = "Test ISSUE for KLoops"
)

func TestTekton(t *testing.T) {
	RunWithReporters(t, "KLoops integration")
}

var _ = CatTests()

func CatTests() bool {
	return Describe("Cat plugin", func() {
		var (
			server  string
			admin   user
			token   *scm.UserToken
			org     *scm.Organization
			orgRepo *scm.Repository
			gitRepo *git.Repo
		)
		BeforeSuite(func() {
			server = GetGitServerURL()
			admin = SetupBasicAuthUser(server, "gitea", "admin")
			token = admin.CreateToken("kloops")
			org = admin.CreateOrganization("kloops-test")
			orgRepo = admin.CreateRepository(org.Name, "kloops-test")
			gitRepo = admin.CloneRepository(server, orgRepo.FullName)
			hmacToken := SetupHmac()
			admin.CreateWebhook(orgRepo, hmacToken)
			InitRepository(gitRepo, "master")
			SetupConfig(orgRepo, ns, hmacToken, token.Token)
		})
		AfterSuite(func() {
			gitRepo.Clean()
			admin.DeleteToken(token.ID)
			admin.DeleteRepository(orgRepo.FullName)
			admin.DeleteOrganization(org.Name)
		})
		Describe("With pull requests", func() {
			var pr *scm.PullRequest
			count := 0
			BeforeEach(func() {
				count++
				branch := fmt.Sprintf("branch-%d", count)
				SetupBranch(gitRepo, "master", branch)
				pr = admin.CreatePullRequest(orgRepo.FullName, branch, "master", prTitle, prBody)
			})
			It("Checks /meow works", func() {
				admin.CreateComment(orgRepo.FullName, "/meow", pr.Number, true)
				WaitForComment(admin.client, orgRepo.FullName, pr.Number, true, `!\[cat image\]`)
			})
			It("Checks /meowvie works", func() {
				admin.CreateComment(orgRepo.FullName, "/meowvie", pr.Number, true)
				WaitForComment(admin.client, orgRepo.FullName, pr.Number, true, `!\[cat image\]`)
			})
		})
		Describe("With issues", func() {
			var issue *scm.Issue
			BeforeEach(func() {
				issue = admin.CreateIssue(orgRepo.FullName, issueTitle, issueBody)
			})
			It("Checks /meow works", func() {
				admin.CreateComment(orgRepo.FullName, "/meow", issue.Number, true)
				WaitForComment(admin.client, orgRepo.FullName, issue.Number, true, `!\[cat image\]`)
			})
			It("Checks /meowvie works", func() {
				admin.CreateComment(orgRepo.FullName, "/meowvie", issue.Number, true)
				WaitForComment(admin.client, orgRepo.FullName, issue.Number, true, `!\[cat image\]`)
			})
		})
	})
}
