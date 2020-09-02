package utils

import (
	"errors"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/git"
)

func GitTokenFunc(rc *v1alpha1.RepoConfig) (git.TokenFunc, error) {
	if rc.Spec.GitHub != nil {
		return func() []byte { return []byte(rc.Spec.GitHub.Token) }, nil
	}
	return nil, errors.New("failed to deduce scm infos from repo config")
}

func GitClone(repoConfig *v1alpha1.RepoConfig, gitClient git.Client, repo string, base string) (*git.Repo, error) {
	tokenFunc, err := GitTokenFunc(repoConfig)
	if err != nil {
		return nil, err
	}
	return gitClient.Clone(repo, base, "toto", tokenFunc)
}
