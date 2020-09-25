package label

import (
	"fmt"
	"strings"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/jenkins-x/go-scm/scm"
)

const pluginName = "label"

var (
	defaultLabels           = []string{"kind", "priority", "area"}
	nonExistentLabelOnIssue = "Those labels are not set on the issue: `%v`"
)

var (
	plugin = plugins.Plugin{
		Description:        "The label plugin provides commands that add or remove certain types of labels. Labels of the following types can be manipulated: 'area/*', 'committee/*', 'kind/*', 'language/*', 'priority/*', 'sig/*', 'triage/*', and 'wg/*'. More labels can be configured to be used via the /label command.",
		ConfigHelpProvider: configHelp,
		Commands: []plugins.Command{{
			Prefix: "remove-",
			Name:   "area|committee|kind|language|priority|sig|triage|wg|label",
			Arg: &plugins.CommandArg{
				Pattern: ".*",
			},
			Description: "Applies or removes a label from one of the recognized types of labels.",
			Action: plugins.
				Invoke(handle).
				When(plugins.Action(scm.ActionCreate)),
		}},
	}
)

func init() {
	plugins.RegisterPlugin(pluginName, plugin)
}

func configString(labels []string) string {
	var formattedLabels []string
	for _, label := range labels {
		formattedLabels = append(formattedLabels, fmt.Sprintf(`"%s/*"`, label))
	}
	return fmt.Sprintf("The label plugin will work on %s and %s labels.", strings.Join(formattedLabels[:len(formattedLabels)-1], ", "), formattedLabels[len(formattedLabels)-1])
}

func configHelp(config *v1alpha1.PluginConfigSpec) (map[string]string, error) {
	labels := []string{}
	labels = append(labels, defaultLabels...)
	labels = append(labels, config.Label.AdditionalLabels...)
	return map[string]string{
			"": configString(labels),
		},
		nil
}

// Get Labels from Regexp matches
func getLabelsFromREMatches(kind string, target string) []string {
	var labels []string
	for _, label := range strings.Split(target, " ") {
		label = strings.ToLower(kind + "/" + strings.TrimSpace(label))
		labels = append(labels, label)
	}
	return labels
}

// getLabelsFromGenericMatches returns label matches with extra labels if those
// have been configured in the plugin config.
func getLabelsFromGenericMatches(label string, additionalLabels []string) []string {
	var labels []string
	for _, l := range additionalLabels {
		if l == label {
			labels = append(labels, label)
		}
	}
	return labels
}

func handle(match plugins.CommandMatch, request plugins.PluginRequest, event plugins.GenericCommentEvent) error {
	// scmClient := request.ScmClient()

	// 	repoLabels, err := spc.GetRepoLabels(org, repo)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	labels, err := spc.GetIssueLabels(org, repo, e.Number, e.IsPR)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	RepoLabelsExisting := map[string]string{}
	// 	for _, l := range repoLabels {
	// 		RepoLabelsExisting[strings.ToLower(l.Name)] = l.Name
	// 	}
	// 	var (
	// 		nonexistent         []string
	// 		noSuchLabelsOnIssue []string
	// 	)

	// 	// Get labels to add and labels to remove from regexp matches
	// 	var lbls []string
	// 	if kind == "label" {
	// 		lbls = append(lbls, getLabelsFromGenericMatches(target, additionalLabels)...)
	// 	} else {
	// 		lbls = append(lbls, getLabelsFromREMatches(kind, target)...)
	// 	}

	// 	for _, lbl := range lbls {
	// 		if remove {
	// 			if !scmprovider.HasLabel(lbl, labels) {
	// 				noSuchLabelsOnIssue = append(noSuchLabelsOnIssue, lbl)
	// 			} else {
	// 				if _, ok := RepoLabelsExisting[lbl]; !ok {
	// 					nonexistent = append(nonexistent, lbl)
	// 				} else {
	// 					if err := spc.RemoveLabel(org, repo, e.Number, lbl, e.IsPR); err != nil {
	// 						log.WithError(err).Errorf("Failed to remove the following label: %s", lbl)
	// 					}
	// 				}
	// 			}
	// 		} else {
	// 			if !scmprovider.HasLabel(lbl, labels) {
	// 				if _, ok := RepoLabelsExisting[lbl]; !ok {
	// 					nonexistent = append(nonexistent, lbl)
	// 				} else {
	// 					if err := spc.AddLabel(org, repo, e.Number, RepoLabelsExisting[lbl], e.IsPR); err != nil {
	// 						log.WithError(err).Errorf("GitHub failed to add the following label: %s", lbl)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}

	// 	//TODO(grodrigues3): Once labels are standardized, make this reply with a comment.
	// 	if len(nonexistent) > 0 {
	// 		log.Infof("Nonexistent labels: %v", nonexistent)
	// 	}

	// 	// Tried to remove Labels that were not present on the Issue
	// 	if len(noSuchLabelsOnIssue) > 0 {
	// 		msg := fmt.Sprintf(nonExistentLabelOnIssue, strings.Join(noSuchLabelsOnIssue, ", "))
	// 		log.Info(msg)
	// 		return spc.CreateComment(org, repo, e.Number, e.IsPR, plugins.FormatResponseRaw(e.Body, e.Link, spc.QuoteAuthorForComment(e.Author.Login), msg))
	// 	}

	return nil
}
