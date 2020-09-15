package events

import (
	"fmt"
	"strconv"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/eddycharly/kloops/pkg/scmprovider"
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
	if event, ok := event.(*scm.PullRequestHook); ok {
		e.handlePullRequest(repoConfig, pluginConfig, scmClient, event)
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
		repoConfig:   repoConfig,
		pluginConfig: pluginConfig,
		gitClient:    e.gitClient,
		client:       e.client,
		logger:       e.logger.WithName(loggerName),
		namespace:    repoConfig.Namespace,
		scmClient:    scmprovider.NewClient(scmClient),
	}
}

func (e *events) getPlugins(repoConfig *v1alpha1.RepoConfig) map[string]*plugins.Plugin {
	p := map[string]*plugins.Plugin{}
	for _, name := range repoConfig.Spec.PluginConfig.Plugins {
		plugin, err := plugins.GetPlugin(name)
		if err != nil {
			e.logger.WithValues("plugin", name).Error(err, "plugin not found")
		} else {
			p[name] = plugin
		}
	}
	return p
}

func (e *events) handleIssueComment(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfigSpec, scmClient *scm.Client, event *scm.IssueCommentHook) {
	e.logger.Info("handle issue comment")
	e.handleGenericComment(repoConfig, pluginConfig, scmClient,
		plugins.GenericCommentEvent{
			GUID:        strconv.Itoa(event.Comment.ID),
			IsPR:        event.Issue.PullRequest,
			Action:      event.Action,
			Body:        event.Comment.Body,
			Link:        event.Comment.Link,
			Number:      event.Issue.Number,
			Repo:        event.Repo,
			Author:      event.Comment.Author,
			IssueAuthor: event.Issue.Author,
			Assignees:   event.Issue.Assignees,
			IssueState:  event.Issue.State,
			IssueBody:   event.Issue.Body,
			IssueLink:   event.Issue.Link,
		},
	)
}

func (e *events) handlePullRequestComment(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfigSpec, scmClient *scm.Client, event *scm.PullRequestCommentHook) {
	e.logger.Info("handle pull request comment")
	e.handleGenericComment(repoConfig, pluginConfig, scmClient,
		plugins.GenericCommentEvent{
			GUID:        strconv.Itoa(event.Comment.ID),
			IsPR:        true,
			Action:      event.Action,
			Body:        event.Comment.Body,
			Link:        event.Comment.Link,
			Number:      event.PullRequest.Number,
			Repo:        event.Repo,
			Author:      event.Comment.Author,
			IssueAuthor: event.PullRequest.Author,
			Assignees:   event.PullRequest.Assignees,
			IssueState:  event.PullRequest.State,
			IssueBody:   event.PullRequest.Body,
			IssueLink:   event.PullRequest.Link,
		},
	)
}

func (e *events) handlePullRequest(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfigSpec, scmClient *scm.Client, event *scm.PullRequestHook) {
	e.logger.Info("handle pull requesr")
	for name, plugin := range e.getPlugins(repoConfig) {
		if plugin.PullRequestHandler != nil {
			if err := plugin.PullRequestHandler(e.makePluginRequest(repoConfig, pluginConfig, scmClient, name), *event); err != nil {
				e.logger.Error(err, "failed to process event")
			}
		}
	}
}

func (e *events) handleGenericComment(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfigSpec, scmClient *scm.Client, event plugins.GenericCommentEvent) {
	e.logger.Info("handle generic comment")
	for name, plugin := range e.getPlugins(repoConfig) {
		if plugin.GenericCommentHandler != nil {
			if err := plugin.GenericCommentHandler(e.makePluginRequest(repoConfig, pluginConfig, scmClient, name), event); err != nil {
				e.logger.Error(err, "failed to process event")
			}
		}
		for _, cmd := range plugin.Commands {
			if err := cmd.InvokeCommandHandler(event, func(handler plugins.CommandEventHandler, event plugins.GenericCommentEvent, match plugins.CommandMatch) error {
				return handler(match, e.makePluginRequest(repoConfig, pluginConfig, scmClient, name), event)
			}); err != nil {
				e.logger.Error(err, "failed to process event")
			}
		}
	}
}
