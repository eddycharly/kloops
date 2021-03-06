package assign

import (
	"fmt"
	"strings"

	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/go-logr/logr"
	"github.com/jenkins-x/go-scm/scm"
)

const pluginName = "assign"

var (
	plugin = plugins.Plugin{
		Description: "The assign plugin assigns or requests reviews from users. Specific users can be assigned with the command '/assign @user1' or have reviews requested of them with the command '/cc @user1'. If no user is specified the commands default to targeting the user who created the command. Assignments and requested reviews can be removed in the same way that they are added by prefixing the commands with 'un'.",
		Commands: []plugins.Command{{
			Prefix: "un",
			Name:   "cc|assign",
			Arg: &plugins.CommandArg{
				Pattern:  `@?"?[-/\w]+"?(?:[ \t]+@?"?[-/\w]+"?)*`,
				Optional: true,
			},
			Description: "Assigns an assignee to the PR or issue or requests a review from the user(s)",
			WhoCanUse:   "Anyone can use the command, but the target user must be an org member, a repo collaborator, or should have previously commented on the issue or PR.",
			Action: plugins.
				Invoke(handleGenericComment).
				When(plugins.Action(scm.ActionCreate)),
		}},
	}
)

func init() {
	plugins.RegisterPlugin(pluginName, plugin)
}

func handleGenericComment(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	err := handle(match, newAssignHandler(event, request))
	if event.IsPR {
		err = combineErrors(err, handle(match, newReviewHandler(event, request)))
	}
	return err
}

func parseLogins(text string) []string {
	var parts []string
	for _, p := range strings.Split(text, " ") {
		t := strings.Trim(p, "@ \"")
		if t == "" {
			continue
		}
		parts = append(parts, t)
	}
	return parts
}

func combineErrors(err1, err2 error) error {
	if err1 != nil && err2 != nil {
		return fmt.Errorf("two errors: 1) %v 2) %v", err1, err2)
	} else if err1 != nil {
		return err1
	} else {
		return err2
	}
}

// handle is the generic handler for the assign plugin. It uses the handler's regexp and affectedLogins
// functions to identify the users to add and/or remove and then passes the appropriate users to the
// handler's add and remove functions. If add fails to add some of the users, a response comment is
// created where the body of the response is generated by the handler's addFailureResponse function.
func handle(match plugins.CommandMatch, h *handler) error {
	if match.Name != h.command {
		return nil
	}
	users := make(map[string]bool)
	if match.Arg == "" {
		users[h.event.Author.Login] = match.Prefix == ""
	} else {
		for _, login := range parseLogins(match.Arg) {
			users[login] = match.Prefix == ""
		}
	}

	var toAdd, toRemove []string
	for login, add := range users {
		if add {
			toAdd = append(toAdd, login)
		} else {
			toRemove = append(toRemove, login)
		}
	}

	if len(toRemove) > 0 {
		// h.log.Printf("Removing %s from %s/%s#%d: %v", h.userType, org, repo, e.Number, toRemove)
		if err := h.remove(h.event.Repo.FullName, h.event.Number, toRemove); err != nil {
			return err
		}
	}
	if len(toAdd) > 0 {
		// h.log.Printf("Adding %s to %s/%s#%d: %v", h.userType, org, repo, e.Number, toAdd)
		if err := h.add(h.event.Repo.FullName, h.event.Number, toAdd); err != nil {
			// if mu, ok := err.(scmprovider.MissingUsers); ok {
			// 	msg := h.addFailureResponse(mu)
			// 	if len(msg) == 0 {
			// 		return nil
			// 	}
			// 	// if err := h.spc.CreateComment(org, repo, e.Number, e.IsPR,
			// 	// 	plugins.FormatResponseRaw(e.Body, e.Link, h.spc.QuoteAuthorForComment(e.Author.Login), msg)); err != nil {
			// 	// 	return fmt.Errorf("comment err: %v", err)
			// 	// }
			// 	return nil
			// }
			return err
		}
	}
	return nil
}

// handler is a struct that contains data about a github event and provides functions to help handle it.
type handler struct {
	// addFailureResponse generates the body of a response comment in the event that the add function fails.
	// addFailureResponse func(mu scmprovider.MissingUsers) string
	// remove is the function that is called on the affected logins for a command prefixed with 'un'.
	remove func(repo string, number int, users []string) error
	// add is the function that is called on the affected logins for a command with no 'un' prefix.
	add func(repo string, number int, users []string) error

	// event is a pointer to the gitprovider.GenericCommentEvent struct that triggered the handler.
	event plugins.GenericCommentEvent
	// command is the name of the command to be executed (assign or cc).
	command string
	// log is a logrus.Entry used to record actions the handler takes.
	log logr.Logger
	// userType is a string that represents the type of users affected by this handler. (e.g. 'assignees')
	userType string
}

func newAssignHandler(event plugins.GenericCommentEvent, request plugins.PluginRequest) *handler {
	// addFailureResponse := func(mu scmprovider.MissingUsers) string {
	// 	return fmt.Sprintf("GitHub didn't allow me to assign the following users: %s.\n\nNote that only [%s members](https://github.com/orgs/%s/people), repo collaborators and people who have commented on this issue/PR can be assigned. Additionally, issues/PRs can only have 10 assignees at the same time.\nFor more information please see [the contributor guide](https://git.k8s.io/community/contributors/guide/#issue-assignment-in-github)", strings.Join(mu.Users, ", "), org, org)
	// }
	scmClient := request.ScmClient()
	return &handler{
		// addFailureResponse: addFailureResponse,
		remove:   scmClient.Issues.Unassign,
		add:      scmClient.Issues.Assign,
		event:    event,
		command:  "assign",
		log:      request.Logger(),
		userType: "assignee(s)",
	}
}

func newReviewHandler(event plugins.GenericCommentEvent, request plugins.PluginRequest) *handler {
	// org := e.Repo.Namespace
	// addFailureResponse := func(mu scmprovider.MissingUsers) string {
	// 	return fmt.Sprintf("GitHub didn't allow me to request PR reviews from the following users: %s.\n\nNote that only [%s members](https://github.com/orgs/%s/people) and repo collaborators can review this PR, and authors cannot review their own PRs.", strings.Join(mu.Users, ", "), org, org)
	// }
	scmClient := request.ScmClient()
	return &handler{
		// addFailureResponse: addFailureResponse,
		remove:   scmClient.PullRequests.UnrequestReview,
		add:      scmClient.PullRequests.RequestReview,
		event:    event,
		command:  "cc",
		log:      request.Logger(),
		userType: "reviewer(s)",
	}
}
