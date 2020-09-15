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

func (e *events) handleGenericComment(repoConfig *v1alpha1.RepoConfig, pluginConfig *v1alpha1.PluginConfigSpec, scmClient *scm.Client, event plugins.GenericCommentEvent) {
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

// // handlePushEvent handles a push event
// func (s *Server) handlePushEvent(l *logrus.Entry, pe *scm.PushHook) {
// 	repo := pe.Repository()
// 	l = l.WithFields(logrus.Fields{
// 		scmprovider.OrgLogField:  repo.Namespace,
// 		scmprovider.RepoLogField: repo.Name,
// 		"ref":                    pe.Ref,
// 		"head":                   pe.After,
// 	})
// 	l.Info("Push event.")
// 	c := 0
// 	for p, h := range s.getPlugins(pe.Repo.Namespace, pe.Repo.Name) {
// 		if h.PushEventHandler != nil {
// 			s.wg.Add(1)
// 			c++
// 			go func(p string, h plugins.PushEventHandler) {
// 				defer s.wg.Done()
// 				agent := plugins.NewAgent(s.ConfigAgent, s.Plugins, s.ClientAgent, s.ServerURL, l.WithField("plugin", p))
// 				if err := h(agent, *pe); err != nil {
// 					agent.Logger.WithError(err).Error("Error handling PushEvent.")
// 				}
// 			}(p, h.PushEventHandler)
// 		}
// 	}
// 	l.WithField("count", strconv.Itoa(c)).Info("number of push handlers")
// }

// func (s *Server) handlePullRequestEvent(l *logrus.Entry, pr *scm.PullRequestHook) {
// 	l = l.WithFields(logrus.Fields{
// 		scmprovider.OrgLogField:  pr.Repo.Namespace,
// 		scmprovider.RepoLogField: pr.Repo.Name,
// 		scmprovider.PrLogField:   pr.PullRequest.Number,
// 		"author":                 pr.PullRequest.Author.Login,
// 		"url":                    pr.PullRequest.Link,
// 	})
// 	action := pr.Action
// 	l.Infof("Pull request %s.", action)
// 	c := 0
// 	repo := pr.PullRequest.Base.Repo
// 	if repo.Name == "" {
// 		repo = pr.Repo
// 	}
// 	for p, h := range s.getPlugins(repo.Namespace, repo.Name) {
// 		if h.PullRequestHandler != nil {
// 			s.wg.Add(1)
// 			c++
// 			go func(p string, h plugins.PullRequestHandler) {
// 				defer s.wg.Done()
// 				agent := plugins.NewAgent(s.ConfigAgent, s.Plugins, s.ClientAgent, s.ServerURL, l.WithField("plugin", p))
// 				agent.InitializeCommentPruner(
// 					pr.Repo.Namespace,
// 					pr.Repo.Name,
// 					pr.PullRequest.Number,
// 				)
// 				if err := h(agent, *pr); err != nil {
// 					agent.Logger.WithError(err).Error("Error handling PullRequestEvent.")
// 				}
// 			}(p, h.PullRequestHandler)
// 		}
// 	}
// 	l.WithField("count", strconv.Itoa(c)).Info("number of PR handlers")

// 	if !actionRelatesToPullRequestComment(action, l) {
// 		return
// 	}
// 	s.handleGenericComment(
// 		l,
// 		&scmprovider.GenericCommentEvent{
// 			GUID:        pr.GUID,
// 			IsPR:        true,
// 			Action:      action,
// 			Body:        pr.PullRequest.Body,
// 			Link:        pr.PullRequest.Link,
// 			Number:      pr.PullRequest.Number,
// 			Repo:        pr.Repo,
// 			Author:      pr.PullRequest.Author,
// 			IssueAuthor: pr.PullRequest.Author,
// 			Assignees:   pr.PullRequest.Assignees,
// 			IssueState:  pr.PullRequest.State,
// 			IssueBody:   pr.PullRequest.Body,
// 			IssueLink:   pr.PullRequest.Link,
// 		},
// 	)
// }

// // handleBranchEvent handles a branch event
// func (s *Server) handleBranchEvent(entry *logrus.Entry, hook *scm.BranchHook) {
// 	// TODO
// }

// // handleReviewEvent handles a PR review event
// func (s *Server) handleReviewEvent(l *logrus.Entry, re scm.ReviewHook) {
// 	l = l.WithFields(logrus.Fields{
// 		scmprovider.OrgLogField:  re.Repo.Namespace,
// 		scmprovider.RepoLogField: re.Repo.Name,
// 		scmprovider.PrLogField:   re.PullRequest.Number,
// 		"review":                 re.Review.ID,
// 		"reviewer":               re.Review.Author.Login,
// 		"url":                    re.Review.Link,
// 	})
// 	l.Infof("Review %s.", re.Action)
// 	for p, h := range s.getPlugins(re.PullRequest.Base.Repo.Namespace, re.PullRequest.Base.Repo.Name) {
// 		if h.ReviewEventHandler != nil {
// 			s.wg.Add(1)
// 			go func(p string, h plugins.ReviewEventHandler) {
// 				defer s.wg.Done()
// 				agent := plugins.NewAgent(s.ConfigAgent, s.Plugins, s.ClientAgent, s.ServerURL, l.WithField("plugin", p))
// 				agent.InitializeCommentPruner(
// 					re.Repo.Namespace,
// 					re.Repo.Name,
// 					re.PullRequest.Number,
// 				)
// 				if err := h(agent, re); err != nil {
// 					agent.Logger.WithError(err).Error("Error handling ReviewEvent.")
// 				}
// 			}(p, h.ReviewEventHandler)
// 		}
// 	}
// 	action := re.Action
// 	if !actionRelatesToPullRequestComment(action, l) {
// 		return
// 	}
// 	s.handleGenericComment(
// 		l,
// 		&scmprovider.GenericCommentEvent{
// 			GUID:        re.GUID,
// 			IsPR:        true,
// 			Action:      action,
// 			Body:        re.Review.Body,
// 			Link:        re.Review.Link,
// 			Number:      re.PullRequest.Number,
// 			Repo:        re.Repo,
// 			Author:      re.Review.Author,
// 			IssueAuthor: re.PullRequest.Author,
// 			Assignees:   re.PullRequest.Assignees,
// 			IssueState:  re.PullRequest.State,
// 			IssueBody:   re.PullRequest.Body,
// 			IssueLink:   re.PullRequest.Link,
// 		},
// 	)
// }

// func actionRelatesToPullRequestComment(action scm.Action, l *logrus.Entry) bool {
// 	switch action {

// 	case scm.ActionCreate, scm.ActionOpen, scm.ActionSubmitted, scm.ActionEdited, scm.ActionDelete, scm.ActionDismissed, scm.ActionUpdate:
// 		return true

// 	case scm.ActionAssigned,
// 		scm.ActionUnassigned,
// 		scm.ActionReviewRequested,
// 		scm.ActionReviewRequestRemoved,
// 		scm.ActionLabel,
// 		scm.ActionUnlabel,
// 		scm.ActionClose,
// 		scm.ActionReopen,
// 		scm.ActionSync:
// 		return false

// 	default:
// 		l.Errorf(failedCommentCoerceFmt, "pull_request", action.String())
// 		return false
// 	}
// }
