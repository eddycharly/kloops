/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shrug

// const (
// 	pluginName = "shrug"
// 	label      = "¯\\_(ツ)_/¯"
// )

// var (
// 	shrugRe   = regexp.MustCompile(`(?mi)^/(?:lh-)?shrug\s*$`)
// 	unshrugRe = regexp.MustCompile(`(?mi)^/(?:lh-)?unshrug\s*$`)
// )

// func init() {
// 	plugins.RegisterHelpProvider(pluginName, helpProvider)
// 	plugins.RegisterIssueCommentHandler(pluginName, handleIssueComment)
// 	plugins.RegisterPullRequestCommentHandler(pluginName, handlePullRequestComment)
// }

// func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
// 	// The Config field is omitted because this plugin is not configurable.
// 	pluginHelp := &pluginhelp.PluginHelp{
// 		Description: label,
// 	}
// 	pluginHelp.AddCommand(pluginhelp.Command{
// 		Usage:       "/[un]shrug",
// 		Description: label,
// 		Featured:    false,
// 		WhoCanUse:   "Anyone, " + label,
// 		Examples:    []string{"/shrug", "/unshrug"},
// 	})
// 	return pluginHelp, nil
// }

// type scmClient interface {
// 	AddLabel(string, int, string) error
// 	RemoveLabel(string, int, string) error
// 	GetLabels(string, int) ([]*scm.Label, error)
// 	CreateComment(string, int, string) error
// }

// type scmTools interface {
// 	QuoteAuthorForComment(string) string
// }

// // func handleGenericComment(pc plugins.Agent, e scmprovider.GenericCommentEvent) error {
// // 	return handle(pc.SCMProviderClient, pc.Logger, &e)
// // }

// func handleIssueComment(request plugins.PluginRequest, event *scm.IssueCommentHook) error {
// 	scmClient := request.ScmClient()
// 	return handle(scmClient.Issues, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.Issue.Number)
// }

// func handlePullRequestComment(request plugins.PluginRequest, event *scm.PullRequestCommentHook) error {
// 	scmClient := request.ScmClient()
// 	return handle(scmClient.PullRequests, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.PullRequest.Number)
// }

// func handle(client scmClient, scmTools scmTools, logger logr.Logger, repo scm.Repository, action scm.Action, comment scm.Comment, number int) error {
// 	if action != scm.ActionCreate {
// 		return nil
// 	}

// 	wantShrug := false
// 	if shrugRe.MatchString(comment.Body) {
// 		wantShrug = true
// 	} else if unshrugRe.MatchString(comment.Body) {
// 		wantShrug = false
// 	} else {
// 		return nil
// 	}

// 	// Only add the label if it doesn't have it yet.
// 	hasShrug := false
// 	issueLabels, err := client.GetLabels(repo.FullName, number)
// 	if err != nil {
// 		logger. /*.WithValues(org, repo, e.Number)*/ Error(err, "Failed to get the labels")
// 	}
// 	for _, candidate := range issueLabels {
// 		if candidate.Name == label {
// 			hasShrug = true
// 			break
// 		}
// 	}
// 	if hasShrug && !wantShrug {
// 		logger.Info("Removing Shrug label.")
// 		resp := "¯\\\\\\_(ツ)\\_/¯"
// 		if err := client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, resp)); err != nil {
// 			return fmt.Errorf("failed to comment on %s/%s#%d: %v", repo.Namespace, repo.Name, number, err)
// 		}
// 		return client.RemoveLabel(repo.FullName, number, label)
// 	} else if !hasShrug && wantShrug {
// 		logger.Info("Adding Shrug label.")
// 		return client.AddLabel(repo.FullName, number, label)
// 	}
// 	return nil
// }
