package utils

import (
	"errors"
	"net/http"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

func ScmInfos(rc *v1alpha1.RepoConfig) (string, string, scm.SecretFunc, error) {
	if rc.Spec.GitHub != nil {
		return "github",
			rc.Spec.GitHub.ServerURL,
			func(scm.Webhook) (string, error) { return rc.Spec.GitHub.HmacToken, nil },
			nil
	}
	return "", "", nil, errors.New("failed to deduce scm infos from repo config")
}

func ScmClient(rc *v1alpha1.RepoConfig) (*scm.Client, scm.SecretFunc, error) {
	driver, serverURL, secretFunc, err := ScmInfos(rc)
	if err != nil {
		return nil, nil, err
	}
	client, err := factory.NewClient(driver, serverURL, "")
	if err != nil {
		return nil, nil, err
	}
	return client, secretFunc, nil
}

func ParseWebhook(r *http.Request, rc *v1alpha1.RepoConfig) (scm.Webhook, error) {
	client, secretFunc, err := ScmClient(rc)
	if err != nil {
		return nil, err
	}
	webHook, err := client.Webhooks.Parse(r, secretFunc)
	if err != nil {
		return nil, err
	}
	return webHook, nil
}
