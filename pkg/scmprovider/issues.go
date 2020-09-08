package scmprovider

import "github.com/jenkins-x/go-scm/scm"

type Issues struct {
	client *scm.IssueService
}

func (s Issues) CreateComment(repo string, number int, comment string) error {
	return nil
}

func (s Issues) AddLabel(repo string, number int, label string) error {
	return nil
}

func (s Issues) RemoveLabel(repo string, number int, label string) error {
	return nil
}

func (s Issues) GetLabels(repo string, number int) ([]*scm.Label, error) {
	return nil, nil
}
