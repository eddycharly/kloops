package size

import (
	"fmt"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/jenkins-x/go-scm/scm"
)

const pluginName = "size"

var defaultSizes = v1alpha1.Size{
	S:   10,
	M:   30,
	L:   100,
	Xl:  500,
	Xxl: 1000,
}

func init() {
	plugins.RegisterPlugin(
		pluginName,
		plugins.Plugin{
			Description:        "The size plugin manages the 'size/*' labels, maintaining the appropriate label on each pull request as it is updated. Generated files identified by the config file '.generated_files' at the repo root are ignored. Labels are applied based on the total number of lines of changes (additions and deletions).",
			ConfigHelpProvider: configHelp,
			PullRequestHandler: handlePullRequest,
		},
	)
}

func configHelp(config *v1alpha1.PluginConfigSpec) (map[string]string, error) {
	sizes := sizesOrDefault(config.Size)
	return map[string]string{
			"": fmt.Sprintf(`The plugin has the following thresholds:<ul>
<li>size/XS:  0-%d</li>
<li>size/S:   %d-%d</li>
<li>size/M:   %d-%d</li>
<li>size/L:   %d-%d</li>
<li>size/XL:  %d-%d</li>
<li>size/XXL: %d+</li>
</ul>`, sizes.S-1, sizes.S, sizes.M-1, sizes.M, sizes.L-1, sizes.L, sizes.Xl-1, sizes.Xl, sizes.Xxl-1, sizes.Xxl),
		},
		nil
}

func handlePullRequest(request plugins.PluginRequest, event scm.PullRequestHook) error {
	if !isPRChanged(event) {
		return nil
	}

	// logger := request.Logger()
	// scmClient := request.ScmClient()

	// gf, err := genfiles.NewGroup(spc, owner, repo, sha)
	// if err != nil {
	// 	// Continue on parse errors, but warn that something is wrong.
	// 	logger.Error(err, "error while parsing .generated_files")
	// }

	// ga, err := gitattributes.NewGroup(func() ([]byte, error) { return spc.GetFile(owner, repo, ".gitattributes", sha) })
	// if err != nil {
	// 	// Continue on parse errors, but warn that something is wrong.
	// 	logger.Error(err, "error while loading .gitattributes")
	// }

	// changes, err := spc.GetPullRequestChanges(owner, repo, num)
	// if err != nil {
	// 	return fmt.Errorf("can not get PR changes for size plugin: %v", err)
	// }

	// newLabel := bucket(count, sizes).label()
	// var hasLabel bool

	// for _, label := range labels {
	// 	if label.Name == newLabel {
	// 		hasLabel = true
	// 		continue
	// 	}

	// 	if strings.HasPrefix(label.Name, labelPrefix) {
	// 		if err := spc.RemoveLabel(owner, repo, num, label.Name, true); err != nil {
	// 			logger.WithValues("label", label.Name).Error(err, "error while removing label")
	// 		}
	// 	}
	// }

	// if hasLabel {
	// 	return nil
	// }

	// if err := spc.AddLabel(owner, repo, num, newLabel, true); err != nil {
	// 	return fmt.Errorf("error adding label to %s/%s PR #%d: %v", owner, repo, num, err)
	// }

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
func isPRChanged(event scm.PullRequestHook) bool {
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
