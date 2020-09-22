package scmprovider

import (
	"context"

	"github.com/jenkins-x/go-scm/scm"
)

type Git struct {
	client scm.GitService
}

// GetRef retruns the ref from repository
func (s Git) GetRef(repo, ref string) (string, error) {
	answer, _, err := s.client.FindRef(context.Background(), repo, ref)
	return answer, err
}

// DeleteRef deletes the ref from repository
func (s Git) DeleteRef(repo, ref string) error {
	_, err := s.client.DeleteRef(context.Background(), repo, ref)
	return err
}
