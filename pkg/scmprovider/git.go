package scmprovider

import (
	"context"

	"github.com/jenkins-x/go-scm/scm"
)

type Git struct {
	client scm.GitService
}

// DeleteRef deletes the ref from repository
func (s Git) DeleteRef(repo, ref string) error {
	ctx := context.Background()
	_, err := s.client.DeleteRef(ctx, repo, ref)
	return err
}
