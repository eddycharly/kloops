package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"emperror.dev/errors"
	"github.com/jenkins-x/go-scm/scm"
)

// AboutThisBotWithoutCommands contains the message that explains how to interact with the bot.
const AboutThisBotWithoutCommands = "Instructions for interacting with me using PR comments are available [here](https://git.k8s.io/community/contributors/guide/pull-requests.md).  If you have questions or suggestions related to my behavior, please file an issue against the [eddycharly/kloops](https://github.com/eddycharly/kloops/issues/new?title=Command%20issue:) repository."

// AboutThisBotCommands contains the message that links to the commands the bot understand.
const AboutThisBotCommands = "I understand the commands that are listed [here](https://go.k8s.io/bot-commands)."

// AboutThisBot contains the text of both AboutThisBotWithoutCommands and AboutThisBotCommands.
const AboutThisBot = AboutThisBotWithoutCommands + " " + AboutThisBotCommands

// FormatResponse nicely formats a response to a generic reason.
func FormatResponse(scmTools ScmTools, to, message, reason string) string {
	format := `@%s:

%s

<details>

%s

%s
</details>`

	return fmt.Sprintf(format, scmTools.QuoteAuthorForComment(to), message, reason, AboutThisBotWithoutCommands)
}

// FormatSimpleResponse formats a response that does not warrant additional explanation in the
// details section.
func FormatSimpleResponse(scmTools ScmTools, to, message string) string {
	format := `@%s:

%s

<details>

%s
</details>`

	return fmt.Sprintf(format, scmTools.QuoteAuthorForComment(to), message, AboutThisBotWithoutCommands)
}

// FormatCommentResponse nicely formats a response to an issue comment.
func FormatCommentResponse(scmTools ScmTools, ic scm.Comment, s string) string {
	return FormatResponseRaw(scmTools, ic.Body, ic.Link, ic.Author.Login, s)
}

// FormatResponseRaw nicely formats a response for one does not have an issue comment
func FormatResponseRaw(scmTools ScmTools, body, bodyURL, login, reply string) string {
	format := `In response to [this](%s):

%s
`
	// Quote the user's comment by prepending ">" to each line.
	var quoted []string
	for _, l := range strings.Split(body, "\n") {
		quoted = append(quoted, ">"+l)
	}
	return FormatResponse(scmTools, login, reply, fmt.Sprintf(format, bodyURL, strings.Join(quoted, "\n")))
}

type ScmTools interface {
	QuoteAuthorForComment(string) string
}

type ScmComments interface {
	CreateComment(context.Context, string, int, *scm.CommentInput) (*scm.Comment, *scm.Response, error)
}

func CreateComment(scmComments ScmComments, scmTools ScmTools, repo scm.Repository, number int, comment scm.Comment, msg string) error {
	commentInput := scm.CommentInput{
		Body: FormatCommentResponse(scmTools, comment, msg),
	}
	_, response, err := scmComments.CreateComment(context.Background(), repo.FullName, number, &commentInput)
	if err != nil {
		var b bytes.Buffer
		_, cperr := io.Copy(&b, response.Body)
		if cperr != nil {
			return errors.Wrapf(cperr, "response: %s", b.String())
		}
		return errors.Wrapf(err, "response: %s", b.String())
	}
	return nil
}
