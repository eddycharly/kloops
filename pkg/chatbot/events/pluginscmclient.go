package events

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"emperror.dev/errors"
	utils "github.com/eddycharly/kloops/pkg/utils/scm"
	"github.com/jenkins-x/go-scm/scm"
)

type pluginScmClient struct {
	client *scm.Client
}

func (c *pluginScmClient) repositoryName(owner string, repo string) string {
	return fmt.Sprintf("%s/%s", owner, repo)
}

func (c *pluginScmClient) GetClient() *scm.Client {
	return c.client
}

func (c *pluginScmClient) CreateComment(owner, repo string, number int, pr bool, comment string) error {
	fullName := c.repositoryName(owner, repo)
	commentInput := scm.CommentInput{
		Body: comment,
	}
	ctx := context.Background()
	if pr {
		_, response, err := c.client.PullRequests.CreateComment(ctx, fullName, number, &commentInput)
		if err != nil {
			var b bytes.Buffer
			_, cperr := io.Copy(&b, response.Body)
			if cperr != nil {
				return errors.Wrapf(cperr, "response: %s", b.String())
			}
			return errors.Wrapf(err, "response: %s", b.String())
		}

	} else {
		_, response, err := c.client.Issues.CreateComment(ctx, fullName, number, &commentInput)
		if err != nil {
			var b bytes.Buffer
			_, cperr := io.Copy(&b, response.Body)
			if cperr != nil {
				return errors.Wrapf(cperr, "reponse: %s", b.String())
			}
			return errors.Wrapf(err, "response: %s", b.String())
		}
	}
	return nil
}

func (c *pluginScmClient) QuoteAuthorForComment(author string) string {
	return utils.QuoteAuthorForComment(c.client, author)
}
