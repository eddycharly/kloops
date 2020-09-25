package utils

import (
	"errors"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/eddycharly/kloops/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GitTokenFunc(client client.Client, repoConfig *v1alpha1.RepoConfig) (git.TokenFunc, error) {
	if repoConfig.Spec.GitHub != nil {
		return func() []byte {
			data, err := utils.GetSecret(client, repoConfig.Namespace, repoConfig.Spec.GitHub.Token)
			if err == nil {
				return data
			}
			return nil
		}, nil
	}
	return nil, errors.New("failed to deduce scm infos from repo config")
}

func GitClone(client client.Client, repoConfig *v1alpha1.RepoConfig, gitClient git.Client, repo string, base string) (*git.Repo, error) {
	tokenFunc, err := GitTokenFunc(client, repoConfig)
	if err != nil {
		return nil, err
	}
	return gitClient.Clone(repo, base, "toto", tokenFunc)
}
