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

// Package wip will label a PR a work-in-progress if the author provides
// a prefix to their pull request title to the same effect. The submit-
// queue will not merge pull requests with the work-in-progress label.
// The label will be removed when the title changes to no longer begin
// with the prefix.
package wip

// const (
// 	// PluginName defines this plugin's registered name.
// 	pluginName = "wip"
// 	label      = "do-not-merge/work-in-progress"
// )

// var (
// 	titleRegex = regexp.MustCompile(`(?i)^\W?WIP\W`)
// )

// func init() {
// 	plugins.RegisterHelpProvider(pluginName, helpProvider)
// 	plugins.RegisterPullRequestHandler(pluginName, handlePullRequest)
// }

// func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
// 	// Only the Description field is specified because this plugin is not triggered with commands and is not configurable.
// 	return &pluginhelp.PluginHelp{
// 			Description: "The wip (Work In Progress) plugin applies the '" + label + "' Label to pull requests whose title starts with 'WIP' or are in the 'draft' stage, and removes it from pull requests when they remove the title prefix or become ready for review. The '" + label + "' Label is typically used to block a pull request from merging while it is still in progress.",
// 		},
// 		nil
// }

// // Strict subset of SCM.Client methods.
// type scmClient interface {
// 	AddLabel(string, int, string) error
// 	RemoveLabel(string, int, string) error
// 	GetLabels(string, int) ([]*scm.Label, error)
// }

// func handlePullRequest(request plugins.PluginRequest, event *scm.PullRequestHook) error {
// 	logger := request.Logger()
// 	scmClient := request.ScmClient()
// 	return handle(scmClient.PullRequests, logger, event.Repository(), event.Action, event.PullRequest)
// }

// func handle(client scmClient, logger logr.Logger, repo scm.Repository, action scm.Action, pr scm.PullRequest) error {
// 	// These are the only actions indicating the PR title may have changed.
// 	if action != scm.ActionOpen &&
// 		action != scm.ActionReopen &&
// 		action != scm.ActionEdited &&
// 		action != scm.ActionUpdate &&
// 		action != scm.ActionReadyForReview {
// 		return nil
// 	}

// 	currentLabels, err := client.GetLabels(repo.FullName, pr.Number)
// 	if err != nil {
// 		return fmt.Errorf("could not get labels for PR %s:%d in WIP plugin: %v", repo.FullName, pr.Number, err)
// 	}
// 	hasLabel := false
// 	for _, l := range currentLabels {
// 		if l.Name == label {
// 			hasLabel = true
// 		}
// 	}

// 	needsLabel := pr.Draft || titleRegex.MatchString(pr.Title)

// 	if needsLabel && !hasLabel {
// 		if err := client.AddLabel(repo.FullName, pr.Number, label); err != nil {
// 			logger.Error(err, "error while adding label")
// 			return err
// 		}
// 	} else if !needsLabel && hasLabel {
// 		if err := client.RemoveLabel(repo.FullName, pr.Number, label); err != nil {
// 			logger.Error(err, "error while removing label")
// 			return err
// 		}
// 	}
// 	return nil
// }
