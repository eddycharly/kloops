package dog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/go-logr/logr"
	"github.com/jenkins-x/go-scm/scm"
)

var (
	match           = regexp.MustCompile(`(?mi)^/(woof|bark)\s*$`)
	fineRegex       = regexp.MustCompile(`(?mi)^/this-is-fine\s*$`)
	notFineRegex    = regexp.MustCompile(`(?mi)^/this-is-not-fine\s*$`)
	unbearableRegex = regexp.MustCompile(`(?mi)^/this-is-unbearable\s*$`)
	filetypes       = regexp.MustCompile(`(?i)\.(jpg|gif|png)$`)
)

const (
	dogURL                = realPack("https://random.dog/woof.json")
	defaultFineImagesRoot = "https://storage.googleapis.com/this-is-fine-images/"
	fineIMG               = "this_is_fine.png"
	notFineIMG            = "this_is_not_fine.png"
	unbearableIMG         = "this_is_unbearable.jpg"
	pluginName            = "dog"
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

type scmClient interface {
	CreateComment(string, int, string) error
}

type scmTools interface {
	ImageTooBig(string) (bool, error)
	QuoteAuthorForComment(string) string
}

type pack interface {
	readDog(scmTools, string) (string, error)
}

type realPack string

var client = http.Client{}

type dogResult struct {
	URL string `json:"url"`
}

// FormatURL will return the GH markdown to show the image for a specific dogURL.
func FormatURL(dogURL string) (string, error) {
	if dogURL == "" {
		return "", errors.New("empty url")
	}
	src, err := url.ParseRequestURI(dogURL)
	if err != nil {
		return "", fmt.Errorf("invalid url %s: %v", dogURL, err)
	}
	return fmt.Sprintf("[![dog image](%s)](%s)", src, src), nil
}

func (u realPack) readDog(scmTools scmTools, dogURL string) (string, error) {
	if dogURL == "" {
		uri := string(u)
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			return "", fmt.Errorf("could not create request %s: %v", uri, err)
		}
		req.Header.Add("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("could not read dog from %s: %v", uri, err)
		}
		defer resp.Body.Close()
		var a dogResult
		if err = json.NewDecoder(resp.Body).Decode(&a); err != nil {
			return "", err
		}
		dogURL = a.URL
	}

	// GitHub doesn't support videos :(
	if !filetypes.MatchString(dogURL) {
		return "", errors.New("unsupported doggo :( unknown filetype: " + dogURL)
	}
	// checking size, GitHub doesn't support big images
	toobig, err := scmTools.ImageTooBig(dogURL)
	if err != nil {
		return "", err
	} else if toobig {
		return "", errors.New("unsupported doggo :( size too big: " + dogURL)
	}
	return FormatURL(dogURL)
}

func handleIssueComment(request plugins.PluginRequest, event *scm.IssueCommentHook) error {
	scmClient := request.ScmClient()
	return handle(scmClient.Issues, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.Issue.Number, dogURL, defaultFineImagesRoot)
}

func handlePullRequestComment(request plugins.PluginRequest, event *scm.PullRequestCommentHook) error {
	scmClient := request.ScmClient()
	return handle(scmClient.PullRequests, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.PullRequest.Number, dogURL, defaultFineImagesRoot)
}

func handle(client scmClient, scmTools scmTools, logger logr.Logger, repo scm.Repository, action scm.Action, comment scm.Comment, number int, p pack, fineImagesRoot string) error {
	// Only consider new comments.
	if action != scm.ActionCreate {
		return nil
	}
	// Make sure they are requesting a dog
	mat := match.FindStringSubmatch(comment.Body)
	url := ""
	if mat == nil {
		// check is this one of the famous.dog
		if fineRegex.FindStringSubmatch(comment.Body) != nil {
			url = fineImagesRoot + fineIMG
		} else if notFineRegex.FindStringSubmatch(comment.Body) != nil {
			url = fineImagesRoot + notFineIMG
		} else if unbearableRegex.FindStringSubmatch(comment.Body) != nil {
			url = fineImagesRoot + unbearableIMG
		}

		if url == "" {
			return nil
		}
	}

	for i := 0; i < 5; i++ {
		resp, err := p.readDog(scmTools, url)
		if err != nil {
			logger.Error(err, "Failed to get dog img")
			continue
		}

		return client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, resp))
	}

	return errors.New("could not find a valid dog image")
}
