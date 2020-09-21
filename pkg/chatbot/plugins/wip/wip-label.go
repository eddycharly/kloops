package wip

import (
	"fmt"
	"regexp"

	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	// PluginName defines this plugin's registered name.
	pluginName = "wip"
	label      = "do-not-merge/work-in-progress"
)

var (
	titleRegex = regexp.MustCompile(`(?i)^\W?WIP\W`)
)

func init() {
	plugins.RegisterPlugin(
		pluginName,
		plugins.Plugin{
			Description:        "The wip (Work In Progress) plugin applies the '" + label + "' Label to pull requests whose title starts with 'WIP' or are in the 'draft' stage, and removes it from pull requests when they remove the title prefix or become ready for review. The '" + label + "' Label is typically used to block a pull request from merging while it is still in progress.",
			PullRequestHandler: handlePullRequest,
		},
	)
}

func handlePullRequest(request plugins.PluginRequest, pe scm.PullRequestHook) error {
	// These are the only actions indicating the PR title may have changed.
	if pe.Action != scm.ActionOpen &&
		pe.Action != scm.ActionReopen &&
		pe.Action != scm.ActionEdited &&
		pe.Action != scm.ActionUpdate &&
		pe.Action != scm.ActionReadyForReview {
		return nil
	}
	return handle(request, pe)
}

func handle(request plugins.PluginRequest, pe scm.PullRequestHook) error {
	scmClient := request.ScmClient()
	logger := request.Logger()

	currentLabels, err := scmClient.PullRequests.GetLabels(pe.PullRequest.Base.Repo.FullName, pe.PullRequest.Number)
	if err != nil {
		return fmt.Errorf("could not get labels for PR %s:%d in WIP plugin: %v", pe.PullRequest.Base.Repo.FullName, pe.PullRequest.Number, err)
	}
	hasLabel := false
	for _, l := range currentLabels {
		if l.Name == label {
			hasLabel = true
		}
	}

	needsLabel := pe.PullRequest.Draft || titleRegex.MatchString(pe.PullRequest.Title)

	if needsLabel && !hasLabel {
		if err := scmClient.PullRequests.AddLabel(pe.PullRequest.Base.Repo.FullName, pe.PullRequest.Number, label); err != nil {
			logger.Error(err, "error while adding Label")
			return err
		}
	} else if !needsLabel && hasLabel {
		if err := scmClient.PullRequests.RemoveLabel(pe.PullRequest.Base.Repo.FullName, pe.PullRequest.Number, label); err != nil {
			logger.Error(err, "error while removing Label")
			return err
		}
	}
	return nil
}
