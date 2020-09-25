package plugins

import (
	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/eddycharly/kloops/pkg/scmprovider"
	"github.com/go-logr/logr"
	"github.com/jenkins-x/go-scm/scm"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginRequest interface {
	RepoConfig() *v1alpha1.RepoConfig
	PluginConfig() *v1alpha1.PluginConfigSpec
	ScmClient() scmprovider.Client
	GitClient() git.Client
	Client() client.Client
	Logger() logr.Logger
}

// GenericCommentEvent is a fake event type that is instantiated for any github event that contains
// comment like content.
// The specific events that are also handled as GenericCommentEvents are:
// - issue_comment events
// - pull_request_review events
// - pull_request_review_comment events
// - pull_request events with action in ["opened", "edited"]
// - issue events with action in ["opened", "edited"]
//
// Issue and PR "closed" events are not coerced to the "deleted" Action and do not trigger
// a GenericCommentEvent because these events don't actually remove the comment content from GH.
type GenericCommentEvent struct {
	IsPR        bool
	Action      scm.Action
	Body        string
	Link        string
	Number      int
	Repo        scm.Repository
	Author      scm.User
	IssueAuthor scm.User
	Assignees   []scm.User
	IssueState  string
	IssueBody   string
	IssueLink   string
	GUID        string
}

// ConfigHelpProvider defines the function type that constructs help about a plugin configuration.
type ConfigHelpProvider func(*v1alpha1.PluginConfigSpec) (map[string]string, error)

// IssueHandler defines the function contract for a scm.Issue handler.
type IssueHandler func(PluginRequest, scm.Issue) error

// PullRequestHandler defines the function contract for a scm.PullRequest handler.
type PullRequestHandler func(PluginRequest, scm.PullRequestHook) error

// StatusEventHandler defines the function contract for a scm.Status handler.
type StatusEventHandler func(PluginRequest, scm.Status) error

// PushEventHandler defines the function contract for a scm.PushHook handler.
type PushEventHandler func(PluginRequest, scm.PushHook) error

// ReviewEventHandler defines the function contract for a ReviewHook handler.
type ReviewEventHandler func(PluginRequest, scm.ReviewHook) error

// GenericCommentHandler defines the function contract for a scm.Comment handler.
type GenericCommentHandler func(PluginRequest, GenericCommentEvent) error

// CommandEventHandler defines the function contract for a command handler.
type CommandEventHandler func(CommandMatch, PluginRequest, GenericCommentEvent) error

// Plugin defines a plugin and its handlers
type Plugin struct {
	Description           string
	ExcludedProviders     sets.String
	ConfigHelpProvider    ConfigHelpProvider
	IssueHandler          IssueHandler
	PullRequestHandler    PullRequestHandler
	PushEventHandler      PushEventHandler
	ReviewEventHandler    ReviewEventHandler
	StatusEventHandler    StatusEventHandler
	GenericCommentHandler GenericCommentHandler
	Commands              []Command
}

// InvokeCommandHandler calls InvokeHandler on all commands
func (plugin Plugin) InvokeCommandHandler(ce GenericCommentEvent, handler func(CommandEventHandler, GenericCommentEvent, CommandMatch) error) error {
	for _, cmd := range plugin.Commands {
		if err := cmd.InvokeCommandHandler(ce, handler); err != nil {
			return err
		}
	}
	return nil
}

// GetHelp returns plugin help
func (plugin Plugin) GetHelp(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
	var err error
	h := &pluginhelp.PluginHelp{
		Description:       plugin.Description,
		Events:            plugin.GetEvents(),
		ExcludedProviders: plugin.ExcludedProviders.List(),
	}
	if plugin.ConfigHelpProvider != nil {
		h.Config, err = plugin.ConfigHelpProvider(config)
	}
	for _, c := range plugin.Commands {
		h.AddCommand(c.GetHelp())
	}
	return h, err
}

// IsProviderExcluded returns true if the given provider is excluded, false otherwise
func (plugin Plugin) IsProviderExcluded(provider string) bool {
	return plugin.ExcludedProviders.Has(provider)
}

func (plugin Plugin) GetEvents() []string {
	var events []string
	if plugin.IssueHandler != nil {
		events = append(events, "issue")
	}
	if plugin.PullRequestHandler != nil {
		events = append(events, "pull_request")
	}
	if plugin.PushEventHandler != nil {
		events = append(events, "push")
	}
	if plugin.ReviewEventHandler != nil {
		events = append(events, "pull_request_review")
	}
	if plugin.StatusEventHandler != nil {
		events = append(events, "status")
	}
	if plugin.GenericCommentHandler != nil {
		events = append(events, "GenericCommentEvent (any event for user text)")
	}
	return events
}
