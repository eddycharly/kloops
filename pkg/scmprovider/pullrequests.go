package scmprovider

import "github.com/jenkins-x/go-scm/scm"

type PullRequests struct {
	client *scm.PullRequestService
}

func (s PullRequests) CreateComment(repo string, number int, comment string) error {
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
