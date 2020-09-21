package pony

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/clients/theponyapi"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	pluginName = "pony"
)

func createPlugin() plugins.Plugin {
	return plugins.Plugin{
		Description: "The pony plugin adds a pony image to an issue or PR in response to the `/pony` command.",
		Commands: []plugins.Command{{
			Name: "pony",
			Arg: &plugins.CommandArg{
				Pattern:  ".+",
				Optional: true,
			},
			Description: "Add a little pony image to the issue or PR. A particular pony can optionally be named for a picture of that specific pony.",
			Action: plugins.
				Invoke(handle).
				When(plugins.Action(scm.ActionCreate)),
		}},
	}
}

func init() {
	plugins.RegisterPlugin(pluginName, createPlugin())
}

func handle(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	logger := request.Logger()
	scmClient := request.ScmClient()
	// Eventually fetch an image
	image, err := fetchImage(scmClient.Tools, match.Arg)
	if err != nil {
		return err
	}
	// Format the response comment
	rspn, err := formatResponse(image)
	if err != nil {
		logger.Error(err, "Failed to format response")
		return err
	}
	// Create comment
	return sendResponse(scmClient, event, rspn)
}

func fetchImage(scmTools scmprovider.Tools, tags string) (string, error) {
	for i := 0; i < 3; i++ {
		image, err := theponyapi.Search(tags)
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
	return fmt.Sprintf("![pony image](%s)", img), nil
}

func sendResponse(scmClient scmprovider.Client, event plugins.GenericCommentEvent, msg string) error {
	if event.IsPR {
		return scmClient.PullRequests.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, msg))
	}
	return scmClient.Issues.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, msg))
}
