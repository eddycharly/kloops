package pony

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/go-logr/logr"
	"github.com/jenkins-x/go-scm/scm"
)

// Only the properties we actually use.
type ponyResult struct {
	Pony ponyResultPony `json:"pony"`
}

type ponyResultPony struct {
	Representations ponyRepresentations `json:"representations"`
}

type ponyRepresentations struct {
	Full  string `json:"full"`
	Small string `json:"small"`
}

const (
	ponyURL    = realHerd("https://theponyapi.com/api/v1/pony/random")
	pluginName = "pony"
	maxPonies  = 5
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

var client = http.Client{}

type scmClient interface {
	CreateComment(string, int, string) error
}

type scmTools interface {
	ImageTooBig(string) (bool, error)
	QuoteAuthorForComment(string) string
}

type herd interface {
	readPony(scmTools, string) (string, error)
}

type realHerd string

func formatURLs(small, full string) string {
	return fmt.Sprintf("[![pony image](%s)](%s)", small, full)
}

func (h realHerd) readPony(scmTools scmTools, tags string) (string, error) {
	uri := string(h) + "?q=" + url.QueryEscape(tags)
	resp, err := client.Get(uri)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("no pony found")
	}
	var a ponyResult
	if err = json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	embedded := a.Pony.Representations.Small
	tooBig, err := scmTools.ImageTooBig(embedded)
	if err != nil {
		return "", fmt.Errorf("couldn't fetch pony for size check: %v", err)
	}
	if tooBig {
		return "", fmt.Errorf("the pony is too big")
	}
	return formatURLs(a.Pony.Representations.Small, a.Pony.Representations.Full), nil
}

func handleIssueComment(request plugins.PluginRequest, event *scm.IssueCommentHook) error {
	scmClient := request.ScmClient()
	return handle(scmClient.Issues, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.Issue.Number, ponyURL)
}

func handlePullRequestComment(request plugins.PluginRequest, event *scm.PullRequestCommentHook) error {
	scmClient := request.ScmClient()
	return handle(scmClient.PullRequests, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.PullRequest.Number, ponyURL)
}

func handle(client scmClient, scmTools scmTools, logger logr.Logger, repo scm.Repository, action scm.Action, comment scm.Comment, number int, p herd) error {
	// Only consider new comments.
	if action != scm.ActionCreate {
		return nil
	}
	// Make sure they are requesting a pony and don't allow requesting more than 'maxPonies' defined.
	mat := match.FindAllStringSubmatch(comment.Body, maxPonies)
	if mat == nil {
		return nil
	}

	var respBuilder strings.Builder
	var tagsSpecified bool
	for _, tag := range mat {
		for i := 0; i < 5; i++ {
			if tag[1] != "" {
				tagsSpecified = true
			}
			resp, err := p.readPony(scmTools, tag[1])
			if err != nil {
				logger.Error(err, "Failed to get a pony")
				continue
			}
			respBuilder.WriteString(resp + "\n")
			break
		}
	}
	if respBuilder.Len() > 0 {
		return client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, respBuilder.String()))
	}

	var msg string
	if tagsSpecified {
		msg = "Couldn't find a pony matching given tag(s)."
	} else {
		msg = "https://theponyapi.com appears to be down"
	}

	if err := client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, msg)); err != nil {
		logger.Error(err, "Failed to leave comment")
	}

	return errors.New("could not find a valid pony image")
}
