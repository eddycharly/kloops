package hold

import (
	"fmt"

	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	pluginName = "hold"
	label      = "do-not-merge/hold"
)

type hasLabelFunc func(label string, issueLabels []*scm.Label) bool

var (
	plugin = plugins.Plugin{
		Description: "The hold plugin allows anyone to add or remove the '" + label + "' Label from a pull request in order to temporarily prevent the PR from merging without withholding approval.",
		Commands: []plugins.Command{{
			Name:        "hold",
			Description: "Adds or removes the `" + label + "` Label which is used to indicate that the PR should not be automatically merged.",
			Arg: &plugins.CommandArg{
				Pattern:  "cancel",
				Optional: true,
			},
			Action: plugins.
				Invoke(handleGenericComment).
				When(plugins.Action(scm.ActionCreate), plugins.IsPR()),
		}},
	}
)

func init() {
	plugins.RegisterPlugin(pluginName, plugin)
}

func hasLabel(labels []*scm.Label, label string) bool {
	for _, candidate := range labels {
		if candidate.Name == label {
			return true
		}
	}
	return false
}

func handleGenericComment(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	return handle(match, request, event)
}

// handle drives the pull request to the desired state. If any user adds
// a /hold directive, we want to add a label if one does not already exist.
// If they add /hold cancel, we want to remove the label if it exists.
func handle(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	needsLabel := match.Arg == "cancel"
	scmClient := request.ScmClient()
	logger := request.Logger()
	issueLabels, err := scmClient.PullRequests.GetLabels(event.Repo.FullName, event.Number)
	if err != nil {
		return fmt.Errorf("failed to get the labels on %s#%d: %v", event.Repo.FullName, event.Number, err)
	}

	hasLabel := hasLabel(issueLabels, label)
	if hasLabel && !needsLabel {
		logger.Info(fmt.Sprintf("Removing %q Label for %s#%d", label, event.Repo.FullName, event.Number))
		return scmClient.PullRequests.RemoveLabel(event.Repo.FullName, event.Number, label)
	} else if !hasLabel && needsLabel {
		logger.Info(fmt.Sprintf("Adding %q Label for %s#%d", label, event.Repo.FullName, event.Number))
		return scmClient.PullRequests.AddLabel(event.Repo.FullName, event.Number, label)
	}
	return nil
}
