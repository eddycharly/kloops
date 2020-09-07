/*
Copyright 2016 The Kubernetes Authors.

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

// Package size contains a Prow plugin which counts the number of lines changed
// in a pull request, buckets this number into a few size classes (S, L, XL, etc),
// and finally labels the pull request with this size.
package size

import (
	"fmt"
	"strings"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/genfiles"
	"github.com/eddycharly/kloops/pkg/gitattributes"
	"github.com/go-logr/logr"
	"github.com/jenkins-x/go-scm/scm"
)

// The sizes are configurable in the `plugins.yaml` config file; the line constants
// in here represent default values used as fallback if none are provided.
const pluginName = "size"

var defaultSizes = v1alpha1.Size{
	S:   10,
	M:   30,
	L:   100,
	Xl:  500,
	Xxl: 1000,
}

func init() {
	plugins.RegisterHelpProvider(pluginName, helpProvider)
	plugins.RegisterPullRequestHandler(pluginName, handlePullRequest)
}

func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
	sizes := sizesOrDefault(config.Size)
	return &pluginhelp.PluginHelp{
			Description: "The size plugin manages the 'size/*' labels, maintaining the appropriate label on each pull request as it is updated. Generated files identified by the config file '.generated_files' at the repo root are ignored. Labels are applied based on the total number of lines of changes (additions and deletions).",
			Config: map[string]string{
				"": fmt.Sprintf(`The plugin has the following thresholds:<ul>
<li>size/XS:  0-%d</li>
<li>size/S:   %d-%d</li>
<li>size/M:   %d-%d</li>
<li>size/L:   %d-%d</li>
<li>size/XL:  %d-%d</li>
<li>size/XXL: %d+</li>
</ul>`, sizes.S-1, sizes.S, sizes.M-1, sizes.M, sizes.L-1, sizes.L, sizes.Xl-1, sizes.Xl, sizes.Xxl-1, sizes.Xxl),
			},
		},
		nil
}

type scmClient interface {
	AddLabel(owner, repo string, number int, label string, pr bool) error
	RemoveLabel(owner, repo string, number int, label string, pr bool) error
	GetIssueLabels(org, repo string, number int, pr bool) ([]*scm.Label, error)
	GetFile(org, repo, filepath, commit string) ([]byte, error)
	GetPullRequestChanges(org, repo string, number int) ([]*scm.Change, error)
}

func handlePullRequest(request plugins.PluginRequest, event *scm.PullRequestHook) error {
	return handlePR(request.ScmClient(), sizesOrDefault(request.PluginConfig().Size), request.Logger(), event)
}

func handlePR(spc scmClient, sizes v1alpha1.Size, logger logr.Logger, event *scm.PullRequestHook) error {
	if !isPRChanged(event) {
		return nil
	}

	var (
		owner = event.PullRequest.Base.Repo.Namespace
		repo  = event.PullRequest.Base.Repo.Name
		num   = event.PullRequest.Number
		sha   = event.PullRequest.Base.Sha
	)

	gf, err := genfiles.NewGroup(spc, owner, repo, sha)
	if err != nil {
		// Continue on parse errors, but warn that something is wrong.
		logger.Error(err, "error while parsing .generated_files")
	}

	ga, err := gitattributes.NewGroup(func() ([]byte, error) { return spc.GetFile(owner, repo, ".gitattributes", sha) })
	if err != nil {
		// Continue on parse errors, but warn that something is wrong.
		logger.Error(err, "error while loading .gitattributes")
	}

	changes, err := spc.GetPullRequestChanges(owner, repo, num)
	if err != nil {
		return fmt.Errorf("can not get PR changes for size plugin: %v", err)
	}

	var count int
	for _, change := range changes {
		// Skip generated and linguist-generated files.
		if gf != nil && ga != nil && (gf.Match(change.Path) || ga.IsLinguistGenerated(change.Path)) {
			continue
		}

		count += change.Additions + change.Deletions
	}

	labels, err := spc.GetIssueLabels(owner, repo, num, true)
	if err != nil {
		logger.Error(err, "while retrieving labels, error")
	}

	newLabel := bucket(count, sizes).label()
	var hasLabel bool

	for _, label := range labels {
		if label.Name == newLabel {
			hasLabel = true
			continue
		}

		if strings.HasPrefix(label.Name, labelPrefix) {
			if err := spc.RemoveLabel(owner, repo, num, label.Name, true); err != nil {
				logger.WithValues("label", label.Name).Error(err, "error while removing label")
			}
		}
	}

	if hasLabel {
		return nil
	}

	if err := spc.AddLabel(owner, repo, num, newLabel, true); err != nil {
		return fmt.Errorf("error adding label to %s/%s PR #%d: %v", owner, repo, num, err)
	}

	return nil
}

// One of a set of discrete buckets.
type size int

const (
	sizeXS size = iota
	sizeS
	sizeM
	sizeL
	sizeXL
	sizeXXL
)

const (
	labelPrefix = "size/"

	labelXS     = "size/XS"
	labelS      = "size/S"
	labelM      = "size/M"
	labelL      = "size/L"
	labelXL     = "size/XL"
	labelXXL    = "size/XXL"
	labelUnkown = "size/?"
)

func (s size) label() string {
	switch s {
	case sizeXS:
		return labelXS
	case sizeS:
		return labelS
	case sizeM:
		return labelM
	case sizeL:
		return labelL
	case sizeXL:
		return labelXL
	case sizeXXL:
		return labelXXL
	}

	return labelUnkown
}

func bucket(lineCount int, sizes v1alpha1.Size) size {
	if lineCount < sizes.S {
		return sizeXS
	} else if lineCount < sizes.M {
		return sizeS
	} else if lineCount < sizes.L {
		return sizeM
	} else if lineCount < sizes.Xl {
		return sizeL
	} else if lineCount < sizes.Xxl {
		return sizeXL
	}

	return sizeXXL
}

// These are the only actions indicating the code diffs may have changed.
func isPRChanged(event *scm.PullRequestHook) bool {
	switch event.Action {
	case scm.ActionOpen:
		return true
	case scm.ActionReopen:
		return true
	case scm.ActionSync:
		return true
	case scm.ActionEdited:
		return true
	case scm.ActionUpdate:
		return true
	default:
		return false
	}
}

func defaultIfZero(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

func sizesOrDefault(sizes v1alpha1.Size) v1alpha1.Size {
	sizes.S = defaultIfZero(sizes.S, defaultSizes.S)
	sizes.M = defaultIfZero(sizes.M, defaultSizes.M)
	sizes.L = defaultIfZero(sizes.L, defaultSizes.L)
	sizes.Xl = defaultIfZero(sizes.Xl, defaultSizes.Xl)
	sizes.Xxl = defaultIfZero(sizes.Xxl, defaultSizes.Xxl)
	return sizes
}
