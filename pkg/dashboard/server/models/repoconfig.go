package models

import (
	"github.com/eddycharly/kloops/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GitHubRepo defines a GitHub repository
type GitHubRepo struct {
	// Owner is the repository owner name
	Owner string `json:"owner"`
	// Repo is the repository owner name
	Repo string `json:"repo"`
	// ServerURL is the GitHub server url
	ServerURL string `json:"server,omitempty"`
	// HmacToken is the secret used to validate webhooks
	HmacToken string `json:"hmacToken,omitempty"`
	// Token is the token used to interact with the git repository
	Token string `json:"token,omitempty"`
}

// GiteaRepo defines a Gitea repository
type GiteaRepo struct {
	// Owner is the repository owner name
	Owner string `json:"owner"`
	// Repo is the repository owner name
	Repo string `json:"repo"`
	// ServerURL is the GitHub server url
	ServerURL string `json:"server,omitempty"`
	// HmacToken is the secret used to validate webhooks
	HmacToken string `json:"hmacToken,omitempty"`
	// Token is the token used to interact with the git repository
	Token string `json:"token,omitempty"`
}

// RepoPluginConfig defines a PluginConfig (it can be a ref or an inline spec)
type RepoPluginConfig struct {
	Ref string `json:"ref,omitempty"`
	// Spec    *PluginConfigSpec `json:"spec,omitempty"`
	Plugins []string `json:"plugins,omitempty"`
}

// RepoConfigSpec defines the desired state of RepoConfig
type RepoConfigSpec struct {
	// BotName is the bot name used by plugins
	BotName string `json:"botName,omitempty"`
	// GitHub defines the GitHub repository details
	GitHub *GitHubRepo `json:"gitHub,omitempty"`
	// Gitea defines the Gitea repository details
	Gitea *GiteaRepo `json:"gitea,omitempty"`
	// AutoMerge configuration for the repository
	AutoMerge *v1alpha1.AutoMerge `json:"autoMerge,omitempty"`
	// PluginConfig defines the plugin configuration for the repository
	PluginConfig RepoPluginConfig `json:"pluginConfig"`
}

// RepoConfig is the Schema for the repoconfigs API
type RepoConfig struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RepoConfigSpec `json:"spec,omitempty"`
}

func ConvertGithubFrom(from *v1alpha1.GitHubRepo) *GitHubRepo {
	if from == nil {
		return nil
	}
	return &GitHubRepo{
		Owner:     from.Owner,
		Repo:      from.Repo,
		ServerURL: from.ServerURL,
	}
}

func ConvertGiteaFrom(from *v1alpha1.GiteaRepo) *GiteaRepo {
	if from == nil {
		return nil
	}
	return &GiteaRepo{
		Owner:     from.Owner,
		Repo:      from.Repo,
		ServerURL: from.ServerURL,
	}
}

func ConvertSpecFrom(from v1alpha1.RepoConfigSpec) RepoConfigSpec {
	return RepoConfigSpec{
		BotName:      from.BotName,
		AutoMerge:    from.AutoMerge,
		PluginConfig: RepoPluginConfig{},
		GitHub:       ConvertGithubFrom(from.GitHub),
		Gitea:        ConvertGiteaFrom(from.Gitea),
	}
}

func ConvertFrom(from v1alpha1.RepoConfig) RepoConfig {
	return RepoConfig{
		ObjectMeta: from.ObjectMeta,
		Spec:       ConvertSpecFrom(from.Spec),
	}
}

func ConvertList(from []v1alpha1.RepoConfig) []RepoConfig {
	var list []RepoConfig
	for _, i := range from {
		list = append(list, ConvertFrom(i))
	}
	return list
}
