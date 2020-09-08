package pony

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/clients/theponyapi"
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
	pluginName = "pony"
)

var (
	match = regexp.MustCompile(`(?mi)^/(?:pony)(?: +(.+?))?\s*$`)
)

func init() {
	plugins.RegisterHelpProvider(pluginName, helpProvider)
	plugins.RegisterIssueCommentHandler(pluginName, handleIssueComment)
	plugins.RegisterPullRequestCommentHandler(pluginName, handlePullRequestComment)
}

func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
	// The Config field is omitted because this plugin is not configurable.
	pluginHelp := &pluginhelp.PluginHelp{
		Description: "The pony plugin adds a pony image to an issue or PR in response to the `/pony` command.",
	}
	pluginHelp.AddCommand(pluginhelp.Command{
		Usage:       "/(pony) [pony]",
		Description: "Add a little pony image to the issue or PR. A particular pony can optionally be named for a picture of that specific pony.",
		Featured:    false,
		WhoCanUse:   "Anyone",
		Examples:    []string{"/pony", "/pony Twilight Sparkle"},
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
	// Make sure they are requesting a cat
	mat := match.FindStringSubmatch(comment.Body)
	if mat == nil {
		return nil
	}
	// Fetch image
	image, err := fetchImage(scmTools, mat[1])
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
	return client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, rspn))
}

func fetchImage(scmTools scmTools, tags string) (string, error) {
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
