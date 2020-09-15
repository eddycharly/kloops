package dog

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/clients/randomdog"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	fineURL       = "https://storage.googleapis.com/this-is-fine-images/this_is_fine.png"
	notFineURL    = "https://storage.googleapis.com/this-is-fine-images/this_is_not_fine.png"
	unbearableURL = "https://storage.googleapis.com/this-is-fine-images/this_is_unbearable.jpg"
	pluginName    = "dog"
)

func createPlugin() plugins.Plugin {
	return plugins.Plugin{
		Description: "The dog plugin adds a dog image to an issue or PR in response to the `/woof` command.",
		Commands: []plugins.Command{{
			Name:        "woof|bark",
			Description: "Add a dog image to the issue or PR",
			Action: plugins.
				Invoke(handle).
				When(plugins.Action(scm.ActionCreate)),
		}, {
			Name:        "this-is-fine",
			Description: "Add a dog image to the issue or PR",
			Action: plugins.
				Invoke(func(_ plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
					return formatURLAndSendResponse(request, event, fineURL)
				}).
				When(plugins.Action(scm.ActionCreate)),
		}, {
			Name:        "this-is-not-fine",
			Description: "Add a dog image to the issue or PR",
			Action: plugins.
				Invoke(func(_ plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
					return formatURLAndSendResponse(request, event, notFineURL)
				}).
				When(plugins.Action(scm.ActionCreate)),
		}, {
			Name:        "this-is-unbearable",
			Description: "Add a dog image to the issue or PR",
			Action: plugins.
				Invoke(func(_ plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
					return formatURLAndSendResponse(request, event, unbearableURL)
				}).
				When(plugins.Action(scm.ActionCreate)),
		}},
	}
}

func init() {
	plugins.RegisterPlugin(pluginName, createPlugin())
}

func formatURLAndSendResponse(request plugins.PluginRequest, event plugins.GenericCommentEvent, url string) error {
	msg, err := formatResponse(url)
	if err != nil {
		return err
	}
	return sendResponse(request, event, msg)
}

func handle(_ plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	logger := request.Logger()
	scmClient := request.ScmClient()
	// Eventually fetch an image
	image, err := fetchImage(scmClient.Tools)
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
	if event.IsPR {
		return scmClient.PullRequests.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, rspn))
	} else {
		return scmClient.Issues.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, rspn))
	}
}

func fetchImage(scmTools scmprovider.Tools) (string, error) {
	for i := 0; i < 5; i++ {
		image, err := randomdog.Get()
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
	return fmt.Sprintf("![dog image](%s)", img), nil
}

func sendResponse(request plugins.PluginRequest, event plugins.GenericCommentEvent, msg string) error {
	scmClient := request.ScmClient()
	if event.IsPR {
		return scmClient.PullRequests.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, msg))
	} else {
		return scmClient.Issues.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, msg))
	}
}
