package cat

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/clients/thecatapi"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/eddycharly/kloops/pkg/utils"
	"github.com/jenkins-x/go-scm/scm"
)

var (
	grumpyKeywords = regexp.MustCompile(`(?mi)^(no|grumpy)\s*$`)
)

const (
	pluginName = "cat"
	grumpyURL  = "https://upload.wikimedia.org/wikipedia/commons/e/ee/Grumpy_Cat_by_Gage_Skidmore.jpg"
)

var (
	plugin = plugins.Plugin{
		Description:        "The cat plugin adds a cat image to an issue or PR in response to the `/meow` command.",
		ConfigHelpProvider: configHelp,
		Commands: []plugins.Command{{
			Name: "meow|meowvie",
			Arg: &plugins.CommandArg{
				Pattern:  `.+`,
				Optional: true,
			},
			Description: "Add a cat image to the issue or PR",
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
			"": fmt.Sprintf("The cat plugin uses an api key for thecatapi.com."),
		},
		nil
}

func handle(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	logger := request.Logger()
	scmClient := request.ScmClient()
	// Fetch image
	image, err := fetchImage(scmClient.Tools, match.Arg, match.Name == "meowvie", getKey(request)())
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
	if event.IsPR {
		return scmClient.PullRequests.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, rspn))
	} else {
		return scmClient.Issues.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, rspn))
	}
}

func getKey(request plugins.PluginRequest) func() string {
	return func() string {
		key, err := utils.GetSecret(request.Client(), request.RepoConfig().Namespace, request.PluginConfig().Cat.Key)
		if err == nil {
			return string(key)
		}
		return ""
	}
}

func fetchImage(scmTools scmprovider.Tools, category string, movieCat bool, key string) (string, error) {
	if grumpyKeywords.MatchString(category) {
		return grumpyURL, nil
	}
	for i := 0; i < 3; i++ {
		image, err := thecatapi.Search(category, movieCat, key)
		if err == nil {
			toobig, err := scmTools.ImageTooBig(image)
			if err == nil && !toobig {
				return image, nil
			}
		}
	}
	return "", errors.New("Failed to find a cat image")
}

func formatResponse(image string) (string, error) {
	img, err := url.Parse(image)
	if err != nil {
		return "", fmt.Errorf("invalid image url %s: %v", image, err)
	}
	return fmt.Sprintf("![cat image](%s)", img), nil
}
