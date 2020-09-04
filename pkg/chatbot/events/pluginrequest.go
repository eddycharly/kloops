package events

import (
	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type pluginRequest struct {
	repoConfig   *v1alpha1.RepoConfigSpec
	pluginConfig *v1alpha1.PluginConfigSpec
	scmClient    *pluginScmClient
	gitClient    git.Client
	client       client.Client
	logger       logr.Logger
	namespace    string
}

func (pr *pluginRequest) RepoConfig() *v1alpha1.RepoConfigSpec {
	return pr.repoConfig
}

func (pr *pluginRequest) PluginConfig() *v1alpha1.PluginConfigSpec {
	return pr.pluginConfig
}

func (pr *pluginRequest) ScmClient() plugins.PluginScmClient {
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

func (pr *pluginRequest) Namespace() string {
	return pr.namespace
}
