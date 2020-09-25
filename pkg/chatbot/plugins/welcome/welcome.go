package welcome

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	pluginName            = "welcome"
	defaultWelcomeMessage = "Welcome @{{.PullRequest.Author.Login}}! It looks like this is your first PR to {{.Repo.FullName}} 🎉"
)

func init() {
	plugins.RegisterPlugin(
		pluginName,
		plugins.Plugin{
			Description:        "The welcome plugin posts a welcoming message when it detects a user's first contribution to a repo.",
			ConfigHelpProvider: configHelp,
			PullRequestHandler: handlePullRequest,
		},
	)
}

func configHelp(config *v1alpha1.PluginConfigSpec) (map[string]string, error) {
	welcomeConfig := map[string]string{
		"": fmt.Sprintf("The welcome plugin is configured to post using following welcome template: %s.", config.Welcome.MessageTemplate),
	}
	return welcomeConfig, nil
}

func handlePullRequest(request plugins.PluginRequest, event scm.PullRequestHook) error {
	scmClient := request.ScmClient()
	pluginConfig := request.PluginConfig()
	// Only consider newly opened PRs
	if event.Action != scm.ActionOpen {
		return nil
	}
	// search for PRs from the author in this repo
	issues, err := scmClient.PullRequests.FindByAuthor(event.Repo.FullName, event.Sender.Login)
	if err != nil {
		return err
	}
	// if there are no results, this is the first! post the welcome comment
	if len(issues) == 0 || (len(issues) == 1 && issues[0].Number == event.PullRequest.Number) {
		// load the template, and run it over the PR info
		welcomeTemplate := pluginConfig.Welcome.MessageTemplate
		if welcomeTemplate == "" {
			welcomeTemplate = defaultWelcomeMessage
		}
		parsedTemplate, err := template.New("welcome").Parse(welcomeTemplate)
		if err != nil {
			return err
		}
		var msgBuffer bytes.Buffer
		err = parsedTemplate.Execute(&msgBuffer, event)
		if err != nil {
			return err
		}
		// actually post the comment
		return scmClient.PullRequests.CreateComment(event.Repo.FullName, event.PullRequest.Number, msgBuffer.String())
	}
	return nil
}
