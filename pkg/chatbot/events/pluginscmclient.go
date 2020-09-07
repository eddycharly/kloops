package events

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"emperror.dev/errors"
	utils "github.com/eddycharly/kloops/pkg/utils/scm"
	"github.com/jenkins-x/go-scm/scm"
	"k8s.io/apimachinery/pkg/util/sets"
)

// NoLabelProviders returns a set of provider names that don't support labels.
func NoLabelProviders() sets.String {
	// "coding" is a placeholder provider name from go-scm that we'll use for testing the comment support for label logic.
	return sets.NewString("stash", "coding")
}

type pluginScmClient struct {
	client  *scm.Client
	botName string
}

// BotName returns the bot name
func (c *pluginScmClient) BotName() string {
	return c.botName
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

// ProviderType returns the type of the underlying SCM provider
func (c *pluginScmClient) ProviderType() string {
	return c.client.Driver.String()
}

// SupportsPRLabels returns true if the underlying provider supports PR labels
func (c *pluginScmClient) SupportsPRLabels() bool {
	return !NoLabelProviders().Has(c.ProviderType())
}

// GetIssueLabels returns the issue labels
func (c *pluginScmClient) GetIssueLabels(org, repo string, number int, pr bool) ([]*scm.Label, error) {
	ctx := context.Background()
	fullName := c.repositoryName(org, repo)
	var allLabels []*scm.Label
	var resp *scm.Response
	var labels []*scm.Label
	var err error
	firstRun := false
	opts := scm.ListOptions{
		Page: 1,
	}
	if pr {
		if !c.SupportsPRLabels() {
			return GetLabelsFromComment(c, org, repo, number)
		}
		for !firstRun || (resp != nil && opts.Page <= resp.Page.Last) {
			labels, resp, err = c.client.PullRequests.ListLabels(ctx, fullName, number, opts)
			if err != nil {
				return nil, err
			}
			firstRun = true
			allLabels = append(allLabels, labels...)
			opts.Page++
		}
	}
	for !firstRun || (resp != nil && opts.Page <= resp.Page.Last) {
		labels, resp, err = c.client.Issues.ListLabels(ctx, fullName, number, opts)
		if err != nil {
			return nil, err
		}
		firstRun = true
		allLabels = append(allLabels, labels...)
		opts.Page++
	}
	return labels, err
}

// AddLabel adds a label
func (c *pluginScmClient) AddLabel(owner, repo string, number int, label string, pr bool) error {
	ctx := context.Background()
	fullName := c.repositoryName(owner, repo)
	if pr {
		if !c.SupportsPRLabels() {
			return AddLabelToComment(c, owner, repo, number, label)
		}
		_, err := c.client.PullRequests.AddLabel(ctx, fullName, number, label)
		return err
	}
	_, err := c.client.Issues.AddLabel(ctx, fullName, number, label)
	return err
}

// RemoveLabel removes labesl
func (c *pluginScmClient) RemoveLabel(owner, repo string, number int, label string, pr bool) error {
	ctx := context.Background()
	fullName := c.repositoryName(owner, repo)
	if pr {
		if !c.SupportsPRLabels() {
			return DeleteLabelFromComment(c, owner, repo, number, label)
		}
		_, err := c.client.PullRequests.DeleteLabel(ctx, fullName, number, label)
		return err
	}
	_, err := c.client.Issues.DeleteLabel(ctx, fullName, number, label)
	return err
}

// EditComment edit a comment
func (c *pluginScmClient) EditComment(owner, repo string, number int, pr bool, id int, comment string) error {
	fullName := c.repositoryName(owner, repo)
	commentInput := scm.CommentInput{
		Body: comment,
	}
	ctx := context.Background()
	if pr {
		_, response, err := c.client.PullRequests.EditComment(ctx, fullName, number, id, &commentInput)
		if err != nil {
			var b bytes.Buffer
			_, cperr := io.Copy(&b, response.Body)
			if cperr != nil {
				return errors.Wrapf(cperr, "response: %s", b.String())
			}
			return errors.Wrapf(err, "response: %s", b.String())
		}

	} else {
		_, response, err := c.client.Issues.EditComment(ctx, fullName, number, id, &commentInput)
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

// DeleteComment delete comments
func (c *pluginScmClient) DeleteComment(org, repo string, number int, pr bool, ID int) error {
	ctx := context.Background()
	fullName := c.repositoryName(org, repo)
	if pr {
		_, err := c.client.PullRequests.DeleteComment(ctx, fullName, number, ID)
		return err
	}
	_, err := c.client.Issues.DeleteComment(ctx, fullName, number, ID)
	return err
}

// ListPullRequestComments list pull request comments
func (c *pluginScmClient) ListPullRequestComments(owner, repo string, number int) ([]*scm.Comment, error) {
	ctx := context.Background()
	fullName := c.repositoryName(owner, repo)
	var allComments []*scm.Comment
	var resp *scm.Response
	var comments []*scm.Comment
	var err error
	firstRun := false
	opts := scm.ListOptions{
		Page: 1,
	}
	for !firstRun || (resp != nil && opts.Page <= resp.Page.Last) {
		comments, resp, err = c.client.PullRequests.ListComments(ctx, fullName, number, opts)
		if err != nil {
			return nil, err
		}
		firstRun = true
		allComments = append(allComments, comments...)
		opts.Page++
	}
	return allComments, nil
}

// GetFile retruns the file from GitHub
func (c *pluginScmClient) GetFile(owner, repo, filepath, commit string) ([]byte, error) {
	ctx := context.Background()
	fullName := c.repositoryName(owner, repo)
	answer, _, err := c.client.Contents.Find(ctx, fullName, filepath, commit)
	var data []byte
	if answer != nil {
		data = answer.Data
	}
	return data, err
}

// GetPullRequestChanges returns the changes in a pull request
func (c *pluginScmClient) GetPullRequestChanges(org, repo string, number int) ([]*scm.Change, error) {
	ctx := context.Background()
	fullName := c.repositoryName(org, repo)
	var allChanges []*scm.Change
	var resp *scm.Response
	var changes []*scm.Change
	var err error
	firstRun := false
	opts := scm.ListOptions{
		Page: 1,
	}
	for !firstRun || (resp != nil && opts.Page <= resp.Page.Last) {
		changes, resp, err = c.client.PullRequests.ListChanges(ctx, fullName, number, opts)
		if err != nil {
			return nil, err
		}
		firstRun = true
		allChanges = append(allChanges, changes...)
		opts.Page++
	}
	return allChanges, nil
}
