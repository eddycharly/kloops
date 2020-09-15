package scmprovider

import (
	"bytes"
	"context"
	"io"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/pkg/errors"
)

type PullRequests struct {
	client scm.PullRequestService
}

func (s PullRequests) CreateComment(repo string, number int, comment string) error {
	commentInput := scm.CommentInput{
		Body: comment,
	}
	ctx := context.Background()
	_, response, err := s.client.CreateComment(ctx, repo, number, &commentInput)
	if err != nil {
		var b bytes.Buffer
		_, cperr := io.Copy(&b, response.Body)
		if cperr != nil {
			return errors.Wrapf(cperr, "reponse: %s", b.String())
		}
		return errors.Wrapf(err, "response: %s", b.String())
	}
	return nil
}

func (s PullRequests) AddLabel(repo string, number int, label string) error {
	return nil
}

func (s PullRequests) RemoveLabel(repo string, number int, label string) error {
	return nil
}

func (s PullRequests) GetLabels(repo string, number int) ([]*scm.Label, error) {
	return nil, nil
}
