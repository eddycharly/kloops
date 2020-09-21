package utils

import (
	"context"
	"errors"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/git"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetSecret(client client.Client, namespace string, secret v1alpha1.Secret) ([]byte, error) {
	if secret.Value != "" {
		return []byte(secret.Value), nil
	}
	nn := types.NamespacedName{
		Namespace: namespace,
		Name:      secret.ValueFrom.SecretKeyRef.Name,
	}
	var s corev1.Secret
	err := client.Get(context.Background(), nn, &s)
	if err != nil {
		return nil, err
	}
	value, ok := s.Data[secret.ValueFrom.SecretKeyRef.Key]
	if !ok {
		return nil, errors.New("key not found in secret")
	}
	return value, nil
}

func GitTokenFunc(client client.Client, repoConfig *v1alpha1.RepoConfig) (git.TokenFunc, error) {
	if repoConfig.Spec.GitHub != nil {
		return func() []byte {
			data, err := GetSecret(client, repoConfig.Namespace, repoConfig.Spec.GitHub.Token)
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
