package utils

import (
	"errors"
	"net/http"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ScmInfos(client client.Client, repoConfig *v1alpha1.RepoConfig) (string, string, string, scm.SecretFunc, error) {
	if repoConfig.Spec.GitHub != nil {
		token, err := getSecret(client, repoConfig.Namespace, repoConfig.Spec.GitHub.Token)
		if err != nil {
			return "", "", "", nil, errors.New("failed to read token")
		}
		return "github",
			repoConfig.Spec.GitHub.ServerURL,
			string(token),
			func(scm.Webhook) (string, error) {
				hmac, err := getSecret(client, repoConfig.Namespace, repoConfig.Spec.GitHub.HmacToken)
				if err != nil {
					return "", err
				}
				return string(hmac), nil
			},
			nil
	}
	return "", "", "", nil, errors.New("failed to deduce scm infos from repo config")
}

func ScmClient(client client.Client, repoConfig *v1alpha1.RepoConfig) (*scm.Client, scm.SecretFunc, error) {
	driver, serverURL, token, secretFunc, err := ScmInfos(client, repoConfig)
	if err != nil {
		return nil, nil, err
	}
	scmClient, err := factory.NewClient(driver, serverURL, token)
	if err != nil {
		return nil, nil, err
	}
	return scmClient, secretFunc, nil
}

func ParseWebhook(client client.Client, request *http.Request, repoConfig *v1alpha1.RepoConfig) (scm.Webhook, *scm.Client, error) {
	scmClient, secretFunc, err := ScmClient(client, repoConfig)
	if err != nil {
		return nil, nil, err
	}
	webHook, err := scmClient.Webhooks.Parse(request, secretFunc)
	if err != nil {
		return nil, nil, err
	}
	return webHook, scmClient, nil
}
