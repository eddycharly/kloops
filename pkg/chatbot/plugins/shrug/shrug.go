package shrug

import (
	"fmt"

	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/jenkins-x/go-scm/scm"
	"k8s.io/test-infra/prow/labels"
)

const (
	pluginName = "shrug"
	label      = "¯\\_(ツ)_/¯"
)

var (
	plugin = plugins.Plugin{
		Description: labels.Shrug,
		Commands: []plugins.Command{{
			Prefix:      "un",
			Name:        "shrug",
			Description: "Adds or removes the " + labels.Shrug + " label",
			Action: plugins.
				Invoke(shrug).
				When(plugins.Action(scm.ActionCreate)),
		}},
	}
)

func init() {
	plugins.RegisterPlugin(pluginName, plugin)
}

func shrug(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	scmClient := request.ScmClient()
	if match.Prefix == "un" {
		return removeLabel(scmClient, event)
	}
	return addLabel(scmClient, event)
}

func hasLabel(scmClient scmprovider.Client, event plugins.GenericCommentEvent) (bool, error) {
	var labels []*scm.Label
	var err error
	if event.IsPR {
		labels, err = scmClient.PullRequests.GetLabels(event.Repo.FullName, event.Number)
	} else {
		labels, err = scmClient.PullRequests.GetLabels(event.Repo.FullName, event.Number)
	}
	if err != nil {
		return false, err
	}
	for _, candidate := range labels {
		if candidate.Name == label {
			return true, nil
		}
	}
	return false, nil
}

func addLabel(scmClient scmprovider.Client, event plugins.GenericCommentEvent) error {
	hasLabel, err := hasLabel(scmClient, event)
	if err != nil || hasLabel {
		return err
	}
	if event.IsPR {
		return scmClient.PullRequests.AddLabel(event.Repo.FullName, event.Number, labels.Shrug)
	}
	return scmClient.Issues.AddLabel(event.Repo.FullName, event.Number, labels.Shrug)
}

func removeLabel(scmClient scmprovider.Client, event plugins.GenericCommentEvent) error {
	hasLabel, err := hasLabel(scmClient, event)
	if err != nil || !hasLabel {
		return err
	}
	resp := "¯\\\\\\_(ツ)\\_/¯"
	if event.IsPR {
		err = scmClient.PullRequests.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, resp))
	} else {
		err = scmClient.Issues.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, resp))
	}
	if err != nil {
		return fmt.Errorf("failed to comment on %s#%d: %v", event.Repo.FullName, event.Number, err)
	}
	if event.IsPR {
		return scmClient.PullRequests.RemoveLabel(event.Repo.FullName, event.Number, labels.Shrug)
	}
	return scmClient.Issues.RemoveLabel(event.Repo.FullName, event.Number, labels.Shrug)
}
