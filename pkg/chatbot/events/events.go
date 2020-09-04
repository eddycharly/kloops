package events

import (
	"fmt"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/go-logr/logr"
	"github.com/jenkins-x/go-scm/scm"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Events interface {
	ProcessWebHook(*v1alpha1.RepoConfig, *v1alpha1.PluginConfigSpec, *scm.Client, scm.Webhook) (string, error)
}

type events struct {
	logger    logr.Logger
	gitClient git.Client
	client    client.Client
}

func NewEvents(client client.Client, gitClient git.Client, logger logr.Logger) Events {
	return &events{
		logger:    logger.WithName("Events"),
		gitClient: gitClient,
		client:    client,
	}
}

func (e *events) ProcessWebHook(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfigSpec, scmClient *scm.Client, event scm.Webhook) (string, error) {
	if _, ok := event.(*scm.PingHook); ok {
		return "pong", nil
	}
	if _, ok := event.(*scm.PushHook); ok {
		// o.server.HandlePushEvent(l, pushHook)
		return "processed push hook", nil
	}
	if _, ok := event.(*scm.PullRequestHook); ok {
		// o.server.HandlePullRequestEvent(l, prHook)
		return "processed PR hook", nil
	}
	if _, ok := event.(*scm.BranchHook); ok {
		// o.server.HandleBranchEvent(l, branchHook)
		return "processed branch hook", nil
	}
	if event, ok := event.(*scm.IssueCommentHook); ok {
		e.handleIssueComment(repoConfig, pluginConfig, scmClient, event)
		return "processed issue comment hook", nil
	}
	if event, ok := event.(*scm.PullRequestCommentHook); ok {
		e.handlePullRequestComment(repoConfig, pluginConfig, scmClient, event)
		return "processed PR comment hook", nil
	}
	if _, ok := event.(*scm.ReviewHook); ok {
		// o.server.HandleReviewEvent(l, *prReviewHook)
		return "processed PR review hook", nil
	}
	e.logger.WithValues("kind", event.Kind()).Info("unknown event kind")
	return fmt.Sprintf("unknown hook %s", event.Kind()), nil
}

func (e *events) makePluginRequest(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfigSpec, scmClient *scm.Client, loggerName string) plugins.PluginRequest {
	return &pluginRequest{
		repoConfig:   &repoConfig.Spec,
		pluginConfig: pluginConfig,
		gitClient:    e.gitClient,
		client:       e.client,
		logger:       e.logger.WithName(loggerName),
		namespace:    repoConfig.Namespace,
		scmClient: &pluginScmClient{
			client: scmClient,
		},
	}
}

func (e *events) handleIssueComment(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfigSpec, scmClient *scm.Client, event *scm.IssueCommentHook) {
	for _, plugin := range repoConfig.Spec.PluginConfig.Plugins {
		if handler := plugins.GetIssueCommentHandler(plugin); handler != nil {
			if err := handler(e.makePluginRequest(repoConfig, pluginConfig, scmClient, plugin), event); err != nil {
				e.logger.Error(err, "failed to process event")
			}

		}
	}
}

func (e *events) handlePullRequestComment(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfigSpec, scmClient *scm.Client, event *scm.PullRequestCommentHook) {
	for _, plugin := range repoConfig.Spec.PluginConfig.Plugins {
		if handler := plugins.GetPullRequestCommentHandler(plugin); handler != nil {
			if err := handler(e.makePluginRequest(repoConfig, pluginConfig, scmClient, plugin), event); err != nil {
				e.logger.Error(err, "failed to process event")
			}
		}
	}
}
