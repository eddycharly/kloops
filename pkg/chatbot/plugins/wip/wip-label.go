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

import (
	"fmt"
	"regexp"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/jenkins-x/go-scm/scm"
)

const (
	// PluginName defines this plugin's registered name.
	pluginName = "wip"
	label      = "do-not-merge/work-in-progress"
)

var (
	titleRegex = regexp.MustCompile(`(?i)^\W?WIP\W`)
)

func init() {
	plugins.RegisterHelpProvider(pluginName, helpProvider)
	plugins.RegisterPullRequestHandler(pluginName, handlePullRequest)
}

func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
	// Only the Description field is specified because this plugin is not triggered with commands and is not configurable.
	return &pluginhelp.PluginHelp{
			Description: "The wip (Work In Progress) plugin applies the '" + label + "' Label to pull requests whose title starts with 'WIP' or are in the 'draft' stage, and removes it from pull requests when they remove the title prefix or become ready for review. The '" + label + "' Label is typically used to block a pull request from merging while it is still in progress.",
		},
		nil
}

// Strict subset of SCM.Client methods.
type scmClient interface {
	GetIssueLabels(org, repo string, number int) ([]scm.Label, error)
	AddLabel(owner, repo string, number int, label string) error
	RemoveLabel(owner, repo string, number int, label string) error
}

func handlePullRequest(request plugins.PluginRequest, event *scm.PullRequestHook) error {
	logger := request.Logger()
	scmClient := request.ScmClient()
	// These are the only actions indicating the PR title may have changed.
	if event.Action != scm.ActionOpen &&
		event.Action != scm.ActionReopen &&
		event.Action != scm.ActionEdited &&
		event.Action != scm.ActionUpdate &&
		event.Action != scm.ActionReadyForReview {
		return nil
	}

	var (
		org    = event.PullRequest.Base.Repo.Namespace
		repo   = event.PullRequest.Base.Repo.Name
		number = event.PullRequest.Number
		title  = event.PullRequest.Title
		draft  = event.PullRequest.Draft
	)

	currentLabels, err := scmClient.GetIssueLabels(org, repo, number, true)
	if err != nil {
		return fmt.Errorf("could not get labels for PR %s/%s:%d in WIP plugin: %v", org, repo, number, err)
	}
	hasLabel := false
	for _, l := range currentLabels {
		if l.Name == label {
			hasLabel = true
		}
	}
	needsLabel := draft || titleRegex.MatchString(title)

	if needsLabel && !hasLabel {
		if err := scmClient.AddLabel(org, repo, number, label, true); err != nil {
			logger.Error(err, "error while adding label")
			return err
		}
	} else if !needsLabel && hasLabel {
		if err := scmClient.RemoveLabel(org, repo, number, label, true); err != nil {
			logger.Error(err, "error while removing label")
			return err
		}
	}
	return nil
}
