package dog

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/clients/randomdog"
	"github.com/go-logr/logr"
	"github.com/jenkins-x/go-scm/scm"
)

type scmClient interface {
	CreateComment(string, int, string) error
}

type scmTools interface {
	ImageTooBig(string) (bool, error)
	QuoteAuthorForComment(string) string
}

const (
	defaultFineImagesRoot = "https://storage.googleapis.com/this-is-fine-images/"
	fineIMG               = "this_is_fine.png"
	notFineIMG            = "this_is_not_fine.png"
	unbearableIMG         = "this_is_unbearable.jpg"
	pluginName            = "dog"
)

var (
	match           = regexp.MustCompile(`(?mi)^/(woof|bark)\s*$`)
	fineRegex       = regexp.MustCompile(`(?mi)^/this-is-fine\s*$`)
	notFineRegex    = regexp.MustCompile(`(?mi)^/this-is-not-fine\s*$`)
	unbearableRegex = regexp.MustCompile(`(?mi)^/this-is-unbearable\s*$`)
)

func init() {
	plugins.RegisterHelpProvider(pluginName, helpProvider)
	plugins.RegisterIssueCommentHandler(pluginName, handleIssueComment)
	plugins.RegisterPullRequestCommentHandler(pluginName, handlePullRequestComment)
}

func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
	// The Config field is omitted because this plugin is not configurable.
	pluginHelp := &pluginhelp.PluginHelp{
		Description: "The dog plugin adds a dog image to an issue or PR in response to the `/woof` command.",
	}
	pluginHelp.AddCommand(pluginhelp.Command{
		Usage:       "/(woof|bark|this-is-{fine|not-fine|unbearable})",
		Description: "Add a dog image to the issue or PR",
		Featured:    false,
		WhoCanUse:   "Anyone",
		Examples:    []string{"/woof", "/bark", "/this-is-{fine|not-fine|unbearable}"},
	})
	return pluginHelp, nil
}

func handleIssueComment(request plugins.PluginRequest, event *scm.IssueCommentHook) error {
	scmClient := request.ScmClient()
	return handle(scmClient.Issues, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.Issue.Number)
}

func handlePullRequestComment(request plugins.PluginRequest, event *scm.PullRequestCommentHook) error {
	scmClient := request.ScmClient()
	return handle(scmClient.PullRequests, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.PullRequest.Number)
}

func handle(client scmClient, scmTools scmTools, logger logr.Logger, repo scm.Repository, action scm.Action, comment scm.Comment, number int) error {
	// Only consider new comments.
	if action != scm.ActionCreate {
		return nil
	}
	// Eventually fetch an image
	image, err := fetchImage(scmTools, comment)
	if err != nil {
		return err
	}
	if image == "" {
		return nil
	}
	// Format the response comment
	rspn, err := formatResponse(image)
	if err != nil {
		logger.Error(err, "Failed to format response")
		return err
	}
	return client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, rspn))
}

func fetchImage(scmTools scmTools, comment scm.Comment) (string, error) {
	mat := match.FindStringSubmatch(comment.Body)
	if mat == nil {
		// check is this one of the famous.dog
		if fineRegex.FindStringSubmatch(comment.Body) != nil {
			return defaultFineImagesRoot + fineIMG, nil
		} else if notFineRegex.FindStringSubmatch(comment.Body) != nil {
			return defaultFineImagesRoot + notFineIMG, nil
		} else if unbearableRegex.FindStringSubmatch(comment.Body) != nil {
			return defaultFineImagesRoot + unbearableIMG, nil
		}
		return "", nil
	}
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
