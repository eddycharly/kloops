package scmprovider

import (
	"context"

	"github.com/jenkins-x/go-scm/scm"
)

type Repositories struct {
	client scm.RepositoryService
}

// ListCollaborators list the collaborators to a repository
func (s Repositories) ListCollaborators(repo string) ([]scm.User, error) {
	ctx := context.Background()
	var allCollabs []scm.User
	var resp *scm.Response
	var collabs []scm.User
	var err error
	firstRun := false
	opts := scm.ListOptions{
		Page: 1,
	}
	for !firstRun || (resp != nil && opts.Page <= resp.Page.Last) {
		collabs, resp, err = s.client.ListCollaborators(ctx, repo, opts)
		if err != nil {
			return nil, err
		}
		firstRun = true
		allCollabs = append(allCollabs, collabs...)
		opts.Page++
	}
	return allCollabs, nil
}
