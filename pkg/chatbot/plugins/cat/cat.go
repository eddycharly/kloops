package cat

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/clients/thecatapi"
	"github.com/eddycharly/kloops/pkg/utils"
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
	pluginName = "cat"
	grumpyURL  = "https://upload.wikimedia.org/wikipedia/commons/e/ee/Grumpy_Cat_by_Gage_Skidmore.jpg"
)

var (
	match          = regexp.MustCompile(`(?mi)^/(?:lh-)?meow(vie)?(?: (.+))?\s*$`)
	grumpyKeywords = regexp.MustCompile(`(?mi)^(no|grumpy)\s*$`)
)

func init() {
	plugins.RegisterHelpProvider(pluginName, helpProvider)
	plugins.RegisterIssueCommentHandler(pluginName, handleIssueComment)
	plugins.RegisterPullRequestCommentHandler(pluginName, handlePullRequestComment)
}

func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
	pluginHelp := &pluginhelp.PluginHelp{
		Description: "The cat plugin adds a cat image to an issue or PR in response to the `/meow` command.",
		Config: map[string]string{
			"": "The cat plugin uses an api key for thecatapi.com stored in the plugin config.",
		},
	}
	pluginHelp.AddCommand(pluginhelp.Command{
		Usage:       "/meow(vie) [CATegory]",
		Description: "Add a cat image to the issue or PR",
		Featured:    false,
		WhoCanUse:   "Anyone",
		Examples:    []string{"/meow", "/meow caturday", "/meowvie clothes", "/lh-meow"},
	})
	return pluginHelp, nil
}

func handleIssueComment(request plugins.PluginRequest, event *scm.IssueCommentHook) error {
	scmClient := request.ScmClient()
	return handle(scmClient.Issues, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.Issue.Number, getKey(request))
}

func handlePullRequestComment(request plugins.PluginRequest, event *scm.PullRequestCommentHook) error {
	scmClient := request.ScmClient()
	return handle(scmClient.PullRequests, scmClient.Tools, request.Logger(), event.Repository(), event.Action, event.Comment, event.PullRequest.Number, getKey(request))
}

func handle(client scmClient, scmTools scmTools, logger logr.Logger, repo scm.Repository, action scm.Action, comment scm.Comment, number int, getKey func() string) error {
	// Only consider new comments.
	if action != scm.ActionCreate {
		return nil
	}
	// Make sure they are requesting a cat
	mat := match.FindStringSubmatch(comment.Body)
	if mat == nil {
		return nil
	}
	category, movieCat, err := parseMatch(mat)
	if err != nil {
		return err
	}
	// Fetch image
	image, err := fetchImage(scmTools, category, movieCat, getKey())
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

func parseMatch(mat []string) (string, bool, error) {
	if len(mat) != 3 {
		err := fmt.Errorf("expected 3 capture groups in regexp match, but got %d", len(mat))
		return "", false, err
	}
	category := strings.TrimSpace(mat[2])
	movieCat := len(mat[1]) > 0 // "vie" suffix is present.
	return category, movieCat, nil
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

func fetchImage(scmTools scmTools, category string, movieCat bool, key string) (string, error) {
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
