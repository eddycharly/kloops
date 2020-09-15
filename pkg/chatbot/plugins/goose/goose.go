/*
Copyright 2019 The Kubernetes Authors.

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

// Package goose adds goose images to an issue or PR in response to a /honk comment
package goose

// type scmClient interface {
// 	CreateComment(string, int, string) error
// }

// type scmTools interface {
// 	ImageTooBig(string) (bool, error)
// 	QuoteAuthorForComment(string) string
// }

// const (
// 	pluginName = "goose"
// )

// var (
// 	match = regexp.MustCompile(`(?mi)^/(honk)\s*$`)
// )

// func init() {
// 	plugins.RegisterHelpProvider(pluginName, helpProvider)
// 	plugins.RegisterIssueCommentHandler(pluginName, handleIssueComment)
// 	plugins.RegisterPullRequestCommentHandler(pluginName, handlePullRequestComment)
// }

// func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
// 	pluginHelp := &pluginhelp.PluginHelp{
// 		Description: "The goose plugin adds a goose image to an issue or PR in response to the `/honk` command.",
// 		Config: map[string]string{
// 			"": "The goose plugin uses an api key for unsplash.com stored in the plugin config",
// 		},
// 	}
// 	pluginHelp.AddCommand(pluginhelp.Command{
// 		Usage:       "/honk",
// 		Description: "Add a goose image to the issue or PR",
// 		Featured:    false,
// 		WhoCanUse:   "Anyone",
// 		Examples:    []string{"/honk"},
// 	})
// 	return pluginHelp, nil
// }

// func handleIssueComment(request plugins.PluginRequest, event *scm.IssueCommentHook) error {
// 	scmClient := request.ScmClient()
// 	return handle(scmClient.Issues, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.Issue.Number, getKey(request))
// }

// func handlePullRequestComment(request plugins.PluginRequest, event *scm.PullRequestCommentHook) error {
// 	scmClient := request.ScmClient()
// 	return handle(scmClient.PullRequests, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.PullRequest.Number, getKey(request))
// }

// func handle(client scmClient, scmTools scmTools, logger logr.Logger, repo scm.Repository, action scm.Action, comment scm.Comment, number int, getKey func() string) error {
// 	// Only consider new comments.
// 	if action != scm.ActionCreate {
// 		return nil
// 	}
// 	// Make sure they are requesting a goose
// 	mat := match.FindStringSubmatch(comment.Body)
// 	if mat == nil {
// 		return nil
// 	}
// 	// Fetch image
// 	image, err := fetchImage(scmTools, getKey())
// 	if err != nil {
// 		logger.Error(err, "Failed to get cat img")
// 		return err
// 	}
// 	// Format the response comment
// 	rspn, err := formatResponse(image)
// 	if err != nil {
// 		logger.Error(err, "Failed to format response")
// 		return err
// 	}
// 	return client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, rspn))
// }

// func getKey(request plugins.PluginRequest) func() string {
// 	return func() string {
// 		key, err := utils.GetSecret(request.Client(), request.RepoConfig().Namespace, request.PluginConfig().Goose.Key)
// 		if err == nil {
// 			return string(key)
// 		}
// 		return ""
// 	}
// }

// func fetchImage(scmTools scmTools, key string) (string, error) {
// 	for i := 0; i < 3; i++ {
// 		image, err := unsplash.Search("goose", key)
// 		if err == nil {
// 			toobig, err := scmTools.ImageTooBig(image)
// 			if err == nil && !toobig {
// 				return image, nil
// 			}
// 		}
// 	}
// 	return "", errors.New("Failed to find a goose image")
// }

// func formatResponse(image string) (string, error) {
// 	img, err := url.Parse(image)
// 	if err != nil {
// 		return "", fmt.Errorf("invalid image url %s: %v", image, err)
// 	}
// 	return fmt.Sprintf("![goose image](%s)", img), nil
// }
