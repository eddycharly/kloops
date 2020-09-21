package help

import (
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	pluginName          = "help"
	labelHelp           = "help wanted"
	labelGoodFirstIssue = "good first issue"
)

var (
	helpGuidelinesURL = "https://git.k8s.io/community/contributors/guide/help-wanted.md"
	helpMsgPruneMatch = "This request has been marked as needing help from a contributor."
	helpMsg           = `
This request has been marked as needing help from a contributor.

Please ensure the request meets the requirements listed [here](` + helpGuidelinesURL + `).

If this request no longer meets these requirements, the label can be removed
by commenting with the ` + "`/remove-help`" + ` command.
`
	goodFirstIssueMsgPruneMatch = "This request has been marked as suitable for new contributors."
	goodFirstIssueMsg           = `
This request has been marked as suitable for new contributors.

Please ensure the request meets the requirements listed [here](` + helpGuidelinesURL + "#good-first-issue" + `).

If this request no longer meets these requirements, the label can be removed
by commenting with the ` + "`/remove-good-first-issue`" + ` command.
`
)

var (
	plugin = plugins.Plugin{
		Description: "The help plugin provides commands that add or remove the '" + labelHelp + "' and the '" + labelGoodFirstIssue + "' labels from issues.",
		Commands: []plugins.Command{{
			Prefix:      "remove-",
			Name:        "help|good-first-issue",
			Description: "Applies or removes the '" + labelHelp + "' and '" + labelGoodFirstIssue + "' labels to an issue.",
			WhoCanUse:   "Anyone can trigger this command on a PR.",
			Action: plugins.
				Invoke(func(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
					// TODO
					// cp, err := pc.CommentPruner()
					// if err != nil {
					// 	return err
					// }
					return handle(match, request, event)
				}).
				When(plugins.Action(scm.ActionCreate), plugins.IsNotPR(), plugins.IssueState("open")),
		}},
	}
)

func init() {
	plugins.RegisterPlugin(pluginName, plugin)
}

func hasLabel(labels []*scm.Label, label string) bool {
	for _, candidate := range labels {
		if candidate.Name == label {
			return true
		}
	}
	return false
}

func handle(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	scmClient := request.ScmClient()
	logger := request.Logger()
	// Determine if the issue has the help and the good-first-issue label
	issueLabels, err := scmClient.Issues.GetLabels(event.Repo.FullName, event.Number)
	if err != nil {
		logger.Error(err, "Failed to get issue labels.")
	}
	hasHelp := hasLabel(issueLabels, labelHelp)
	hasGoodFirstIssue := hasLabel(issueLabels, labelGoodFirstIssue)

	// If PR has help label and we're asking for it to be removed, remove label
	if hasHelp && match.Name == "help" && match.Prefix == "remove-" {
		if err := scmClient.Issues.RemoveLabel(event.Repo.FullName, event.Number, labelHelp); err != nil {
			logger.Error(err, "GitHub failed to remove the following label")
		}

		// botName, err := spc.BotName()
		// if err != nil {
		// 	log.WithError(err).Errorf("Failed to get bot name.")
		// }
		// cp.PruneComments(e.IsPR, shouldPrune(log, botName, helpMsgPruneMatch))

		// if it has the good-first-issue label, remove it too
		if hasGoodFirstIssue {
			if err := scmClient.Issues.RemoveLabel(event.Repo.FullName, event.Number, labelGoodFirstIssue); err != nil {
				logger.Error(err, "GitHub failed to remove the following label")
			}
			// cp.PruneComments(e.IsPR, shouldPrune(log, botName, goodFirstIssueMsgPruneMatch))
		}

		return nil
	}

	// If PR does not have the good-first-issue label and we are asking for it to be added,
	// add both the good-first-issue and help labels
	if !hasGoodFirstIssue && match.Name == "good-first-issue" && match.Prefix != "remove-" {
		// if err := spc.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(e.Body, e.IssueLink, spc.QuoteAuthorForComment(commentAuthor), goodFirstIssueMsg)); err != nil {
		// 	log.WithError(err).Errorf("Failed to create comment \"%s\".", goodFirstIssueMsg)
		// }

		if err := scmClient.Issues.AddLabel(event.Repo.FullName, event.Number, labelGoodFirstIssue); err != nil {
			logger.Error(err, "GitHub failed to add the following label")
		}

		if !hasHelp {
			if err := scmClient.Issues.AddLabel(event.Repo.FullName, event.Number, labelHelp); err != nil {
				logger.Error(err, "GitHub failed to add the following label")
			}
		}

		return nil
	}

	// If PR does not have the help label and we're asking it to be added,
	// add the label
	if !hasHelp && match.Name == "help" && match.Prefix != "remove-" {
		// if err := scmClient.Issues.CreateComment(event.Repo.FullName, event.Number, plugins.FormatResponseRaw(e.Body, e.IssueLink, spc.QuoteAuthorForComment(commentAuthor), helpMsg)); err != nil {
		// 	log.WithError(err).Errorf("Failed to create comment \"%s\".", helpMsg)
		// }
		if err := scmClient.Issues.AddLabel(event.Repo.FullName, event.Number, labelHelp); err != nil {
			logger.Error(err, "GitHub failed to add the following label")
		}

		return nil
	}

	// If PR has good-first-issue label and we are asking for it to be removed,
	// remove just the good-first-issue label
	if hasGoodFirstIssue && match.Name == "good-first-issue" && match.Prefix == "remove-" {
		if err := scmClient.Issues.RemoveLabel(event.Repo.FullName, event.Number, labelGoodFirstIssue); err != nil {
			logger.Error(err, "GitHub failed to remove the following label")
		}

		// botName, err := spc.BotName()
		// if err != nil {
		// 	log.WithError(err).Errorf("Failed to get bot name.")
		// }
		// cp.PruneComments(e.IsPR, shouldPrune(log, botName, goodFirstIssueMsgPruneMatch))

		return nil
	}

	return nil
}

// // shouldPrune finds comments left by this plugin.
// func shouldPrune(log *logrus.Entry, botName, msgPruneMatch string) func(*scm.Comment) bool {
// 	return func(comment *scm.Comment) bool {
// 		if comment.Author.Login != botName {
// 			return false
// 		}
// 		return strings.Contains(comment.Body, msgPruneMatch)
// 	}
// }
