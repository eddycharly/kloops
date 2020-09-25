package events

import (
	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type pluginRequest struct {
	repoConfig   *v1alpha1.RepoConfig
	pluginConfig *v1alpha1.PluginConfigSpec
	scmClient    scmprovider.Client
	gitClient    git.Client
	client       client.Client
	logger       logr.Logger
	namespace    string
}

func (pr *pluginRequest) RepoConfig() *v1alpha1.RepoConfig {
	return pr.repoConfig
}

func (pr *pluginRequest) PluginConfig() *v1alpha1.PluginConfigSpec {
	return pr.pluginConfig
}

func (pr *pluginRequest) ScmClient() scmprovider.Client {
	return pr.scmClient
}

func (pr *pluginRequest) GitClient() git.Client {
	return pr.gitClient
}

func (pr *pluginRequest) Client() client.Client {
	return pr.client
}

func (pr *pluginRequest) Logger() logr.Logger {
	return pr.logger
}
