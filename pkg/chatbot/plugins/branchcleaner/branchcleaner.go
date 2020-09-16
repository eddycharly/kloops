package branchcleaner

import (
	"fmt"

	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	pluginName = "branchcleaner"
)

func init() {
	plugins.RegisterPlugin(
		pluginName,
		plugins.Plugin{
			Description:        "The branchcleaner plugin automatically deletes source branches for merged PRs between two branches on the same repository. This is helpful to keep repos that don't allow forking clean.",
			PullRequestHandler: handlePullRequest,
		},
	)
}

func handlePullRequest(request plugins.PluginRequest, event scm.PullRequestHook) error {
	scmClient := request.ScmClient()
	// Only consider closed PRs that got merged
	if event.Action != scm.ActionClose || !event.PullRequest.Merged {
		return nil
	}
	pr := event.PullRequest
	// Only consider PRs from the same repo
	if pr.Base.Repo.FullName != pr.Head.Repo.FullName {
		return nil
	}
	// Delete branch
	branch := scmClient.Tools.RefToBranchName(pr.Head.Ref)
	if err := request.ScmClient().Git.DeleteRef(pr.Base.Repo.FullName, branch); err != nil {
		return fmt.Errorf("failed to delete branch %s on repo %s after Pull Request #%d got merged: %v", branch, pr.Base.Repo.FullName, event.PullRequest.Number, err)
	}
	return nil
}
