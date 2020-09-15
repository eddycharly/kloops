package yuks

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/clients/icanhazdadjoke"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	pluginName = "yuks"
)

var (
	simple = regexp.MustCompile(`^[\w?'!., ]+$`)
)

func createPlugin() plugins.Plugin {
	return plugins.Plugin{
		Description: "The yuks plugin comments with jokes in response to the `/joke` command.",
		Commands: []plugins.Command{{
			Name:        "joke",
			Description: "Tells a joke.",
			Action: plugins.
				Invoke(joke).
				When(plugins.Action(scm.ActionCreate)),
		}},
	}
}

func init() {
	plugins.RegisterPlugin(pluginName, createPlugin())
}

func joke(_ plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	logger := request.Logger()
	scmClient := request.ScmClient()
	// Fetch joke
	joke, err := fetchJoke()
	if err != nil {
		logger.Error(err, "Failed to get joke")
		return err
	}
	// Format the response comment
	rspn := escapeMarkdown(joke)
	// Create comment
	return sendResponse(scmClient, event, rspn)
}

func fetchJoke() (string, error) {
	for i := 0; i < 3; i++ {
		joke, err := icanhazdadjoke.Get()
		if err == nil {
			return joke, nil
		}
	}
	return "", errors.New("Failed to find a joke")
}

// escapeMarkdown takes a string and returns a serialized version of it such that all the symbols
// are treated as text instead of Markdown syntax. It escapes the symbols using numeric character
// references with the decimal notation. See https://www.w3.org/TR/html401/charset.html#h-5.3.1
func escapeMarkdown(s string) string {
	var b bytes.Buffer
	for _, r := range []rune(s) {
		// Check for simple characters as they are considered safe, otherwise we escape the rune.
		c := string(r)
		if simple.MatchString(c) {
			b.WriteString(c)
		} else {
			b.WriteString(fmt.Sprintf("&#%d;", r))
		}
	}
	return b.String()
}

func sendResponse(scmClient scmprovider.Client, event plugins.GenericCommentEvent, msg string) error {
	if event.IsPR {
		return scmClient.PullRequests.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, msg))
	}
	return scmClient.Issues.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(scmClient.Tools, event.Body, event.Link, event.Author.Login, msg))
}
