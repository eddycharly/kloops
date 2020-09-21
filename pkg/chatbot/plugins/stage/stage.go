package stage

import (
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/jenkins-x/go-scm/scm"
)

var (
	stageAlpha  = "stage/alpha"
	stageBeta   = "stage/beta"
	stageStable = "stage/stable"
	stageLabels = []string{stageAlpha, stageBeta, stageStable}
)

const pluginName = "stage"

var (
	plugin = plugins.Plugin{
		Description: "Label the stage of an issue as alpha/beta/stable",
		Commands: []plugins.Command{{
			Name: "stage",
			Arg: &plugins.CommandArg{
				Pattern: "alpha|beta|stable",
			},
			Description: "Labels the stage of an issue as alpha/beta/stable",
			Action: plugins.
				Invoke(stage).
				When(plugins.Action(scm.ActionCreate)),
		}, {
			Name: "remove-stage",
			Arg: &plugins.CommandArg{
				Pattern: "alpha|beta|stable",
			},
			Description: "Removes the stage label of an issue as alpha/beta/stable",
			Action: plugins.
				Invoke(unstage).
				When(plugins.Action(scm.ActionCreate)),
		}},
	}
)

func init() {
	plugins.RegisterPlugin(pluginName, plugin)
}

func unstage(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	lbl := "stage/" + match.Arg
	scmClient := request.ScmClient()
	if event.IsPR {
		return scmClient.PullRequests.RemoveLabel(event.Repo.FullName, event.Number, lbl)
	}
	return scmClient.Issues.RemoveLabel(event.Repo.FullName, event.Number, lbl)
}

func stage(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	lbl := "stage/" + match.Arg
	scmClient := request.ScmClient()
	logger := request.Logger()
	var err error
	for _, label := range stageLabels {
		if event.IsPR {
			err = scmClient.PullRequests.RemoveLabel(event.Repo.FullName, event.Number, label)
		} else {
			err = scmClient.Issues.RemoveLabel(event.Repo.FullName, event.Number, label)
		}
		if err != nil {
			logger.Error(err, "Failed to remove label")
		}
	}
	if event.IsPR {
		err = scmClient.PullRequests.AddLabel(event.Repo.FullName, event.Number, lbl)
	} else {
		err = scmClient.Issues.AddLabel(event.Repo.FullName, event.Number, lbl)
	}
	if err != nil {
		logger.Error(err, "Failed to add label")
	}
	return err
}
