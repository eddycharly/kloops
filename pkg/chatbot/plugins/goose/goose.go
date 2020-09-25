package goose

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/clients/unsplash"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/eddycharly/kloops/pkg/utils"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	pluginName = "goose"
)

var (
	plugin = plugins.Plugin{
		Description:        "The goose plugin adds a goose image to an issue or PR in response to the `/honk` command.",
		ConfigHelpProvider: configHelp,
		Commands: []plugins.Command{{
			Name:        "honk",
			Description: "Add a goose image to the issue or PR",
			Action: plugins.
				Invoke(handle).
				When(plugins.Action(scm.ActionCreate)),
		}},
	}
)

func init() {
	plugins.RegisterPlugin(pluginName, plugin)
}

func configHelp(config *v1alpha1.PluginConfigSpec) (map[string]string, error) {
	return map[string]string{
			"": "The goose plugin uses an api key for unsplash.com stored in the plugin config",
		},
		nil
}

func handle(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	logger := request.Logger()
	scmClient := request.ScmClient()
	// Fetch image
	image, err := fetchImage(scmClient.Tools, getKey(request)())
	if err != nil {
		logger.Error(err, "Failed to get cat img")
		return err
	}
	// Format the response comment
	rspn, err := formatResponse(image)
	if err != nil {
		logger.Error(err, "Failed to format response")
		return err
	}
	// Create comment
	if event.IsPR {
		return scmClient.PullRequests.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, rspn))
	}
	return scmClient.Issues.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, rspn))
}

func getKey(request plugins.PluginRequest) func() string {
	return func() string {
		pluginConfig := request.PluginConfig()
		if pluginConfig == nil {
			return ""
		}
		key, err := utils.GetSecret(request.Client(), request.RepoConfig().Namespace, pluginConfig.Goose.Key)
		if err == nil {
			return string(key)
		}
		fmt.Println(err)
		return ""
	}
}

func fetchImage(scmTools scmprovider.Tools, key string) (string, error) {
	for i := 0; i < 3; i++ {
		image, err := unsplash.Search("goose", key)
		if err == nil {
			toobig, err := scmTools.ImageTooBig(image)
			if err == nil && !toobig {
				return image, nil
			}
		}
	}
	return "", errors.New("Failed to find a goose image")
}

func formatResponse(image string) (string, error) {
	img, err := url.Parse(image)
	if err != nil {
		return "", fmt.Errorf("invalid image url %s: %v", image, err)
	}
	return fmt.Sprintf("![goose image](%s)", img), nil
}
