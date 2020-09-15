package scmprovider

import (
	"bytes"
	"context"
	"io"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/pkg/errors"
)

type Issues struct {
	client scm.IssueService
}

func (s Issues) CreateComment(repo string, number int, comment string) error {
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

func (s Issues) AddLabel(repo string, number int, label string) error {
	ctx := context.Background()
	_, err := s.client.AddLabel(ctx, repo, number, label)
	return err
}

func (s Issues) RemoveLabel(repo string, number int, label string) error {
	ctx := context.Background()
	_, err := s.client.DeleteLabel(ctx, repo, number, label)
	return err
}

func (s Issues) GetLabels(repo string, number int) ([]*scm.Label, error) {
	ctx := context.Background()
	var allLabels []*scm.Label
	var resp *scm.Response
	var labels []*scm.Label
	var err error
	firstRun := false
	opts := scm.ListOptions{
		Page: 1,
	}
	for !firstRun || (resp != nil && opts.Page <= resp.Page.Last) {
		labels, resp, err = s.client.ListLabels(ctx, repo, number, opts)
		if err != nil {
			return nil, err
		}
		firstRun = true
		allLabels = append(allLabels, labels...)
		opts.Page++
	}
	return labels, err
}
