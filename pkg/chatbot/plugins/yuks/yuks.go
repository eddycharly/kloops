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

package yuks

// type scmClient interface {
// 	CreateComment(string, int, string) error
// }

// type scmTools interface {
// 	QuoteAuthorForComment(string) string
// }

// const (
// 	pluginName = "yuks"
// )

// var (
// 	match  = regexp.MustCompile(`(?mi)^/joke\s*$`)
// 	simple = regexp.MustCompile(`^[\w?'!., ]+$`)
// )

// func init() {
// 	plugins.RegisterHelpProvider(pluginName, helpProvider)
// 	plugins.RegisterIssueCommentHandler(pluginName, handleIssueComment)
// 	plugins.RegisterPullRequestCommentHandler(pluginName, handlePullRequestComment)
// }

// func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
// 	// The Config field is omitted because this plugin is not configurable.
// 	pluginHelp := &pluginhelp.PluginHelp{
// 		Description: "The yuks plugin comments with jokes in response to the `/joke` command.",
// 	}
// 	pluginHelp.AddCommand(pluginhelp.Command{
// 		Usage:       "/joke",
// 		Description: "Tells a joke.",
// 		Featured:    false,
// 		WhoCanUse:   "Anyone can use the `/joke` command.",
// 		Examples:    []string{"/joke"},
// 	})
// 	return pluginHelp, nil
// }

// func handleIssueComment(request plugins.PluginRequest, event *scm.IssueCommentHook) error {
// 	scmClient := request.ScmClient()
// 	return handle(scmClient.Issues, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.Issue.Number)
// }

// func handlePullRequestComment(request plugins.PluginRequest, event *scm.PullRequestCommentHook) error {
// 	scmClient := request.ScmClient()
// 	return handle(scmClient.PullRequests, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.PullRequest.Number)
// }

// func handle(client scmClient, scmTools scmTools, logger logr.Logger, repo scm.Repository, action scm.Action, comment scm.Comment, number int) error {
// 	// Only consider new comments.
// 	if action != scm.ActionCreate {
// 		return nil
// 	}
// 	// Make sure they are requesting a joke
// 	if !match.MatchString(comment.Body) {
// 		return nil
// 	}
// 	// Fetch joke
// 	joke, err := fetchJoke()
// 	if err != nil {
// 		logger.Error(err, "Failed to get joke")
// 		return err
// 	}
// 	// Format the response comment
// 	rspn := escapeMarkdown(joke)
// 	return client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, rspn))
// }

// func fetchJoke() (string, error) {
// 	for i := 0; i < 3; i++ {
// 		joke, err := icanhazdadjoke.Get()
// 		if err == nil {
// 			return joke, nil
// 		}
// 	}
// 	return "", errors.New("Failed to find a joke")
// }

// // escapeMarkdown takes a string and returns a serialized version of it such that all the symbols
// // are treated as text instead of Markdown syntax. It escapes the symbols using numeric character
// // references with the decimal notation. See https://www.w3.org/TR/html401/charset.html#h-5.3.1
// func escapeMarkdown(s string) string {
// 	var b bytes.Buffer
// 	for _, r := range []rune(s) {
// 		// Check for simple characters as they are considered safe, otherwise we escape the rune.
// 		c := string(r)
// 		if simple.MatchString(c) {
// 			b.WriteString(c)
// 		} else {
// 			b.WriteString(fmt.Sprintf("&#%d;", r))
// 		}
// 	}
// 	return b.String()
// }
