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

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/eddycharly/kloops/pkg/utils"
	"github.com/go-logr/logr"
	"github.com/jenkins-x/go-scm/scm"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	match = regexp.MustCompile(`(?mi)^/(honk)\s*$`)
	honk  = &realGaggle{
		url: "https://api.unsplash.com/photos/random?query=goose",
	}
)

const (
	pluginName = "goose"
)

func init() {
	plugins.RegisterHelpProvider(pluginName, helpProvider)
	plugins.RegisterIssueCommentHandler(pluginName, handleIssueComment)
	plugins.RegisterPullRequestCommentHandler(pluginName, handlePullRequestComment)
}

func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
	pluginHelp := &pluginhelp.PluginHelp{
		Description: "The goose plugin adds a goose image to an issue or PR in response to the `/honk` command.",
		Config: map[string]string{
			"": "The goose plugin uses an api key for unsplash.com stored in the plugin config",
		},
	}
	pluginHelp.AddCommand(pluginhelp.Command{
		Usage:       "/honk",
		Description: "Add a goose image to the issue or PR",
		Featured:    false,
		WhoCanUse:   "Anyone",
		Examples:    []string{"/honk"},
	})
	return pluginHelp, nil
}

type scmClient interface {
	CreateComment(string, int, string) error
}

type scmTools interface {
	ImageTooBig(string) (bool, error)
	QuoteAuthorForComment(string) string
}

type gaggle interface {
	readGoose(scmTools) (string, error)
}

type realGaggle struct {
	url    string
	lock   sync.RWMutex
	update time.Time
	key    string
}

func (g *realGaggle) setKey(client client.Client, namespace string, secret v1alpha1.Secret) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if !time.Now().After(g.update) {
		return
	}
	g.update = time.Now().Add(1 * time.Minute)
	key, err := utils.GetSecret(client, namespace, secret)
	if err == nil {
		g.key = strings.TrimSpace(string(key))
		return
	}
	// log.WithValues("keyPath", keyPath).Error(err, "failed to read key")
	g.key = ""
}

type gooseResult struct {
	ID     string   `json:"id"`
	Images imageSet `json:"urls"`
}

type imageSet struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

func (gr gooseResult) Format() (string, error) {
	if gr.Images.Small == "" {
		return "", errors.New("empty image url")
	}
	img, err := url.Parse(gr.Images.Small)
	if err != nil {
		return "", fmt.Errorf("invalid image url %s: %v", gr.Images.Small, err)
	}

	return fmt.Sprintf("\n![goose image](%s)", img), nil
}

func (g *realGaggle) URL() string {
	g.lock.RLock()
	defer g.lock.RUnlock()
	uri := string(g.url)
	if g.key != "" {
		uri += "&client_id=" + url.QueryEscape(g.key)
	}
	return uri
}

func (g *realGaggle) readGoose(scmTools scmTools) (string, error) {
	geese := make([]gooseResult, 1)
	uri := g.URL()
	resp, err := http.Get(uri)
	if err != nil {
		return "", fmt.Errorf("could not read goose from %s: %v", uri, err)
	}
	defer resp.Body.Close()
	if sc := resp.StatusCode; sc > 299 || sc < 200 {
		return "", fmt.Errorf("failing %d response from %s", sc, uri)
	}
	if err = json.NewDecoder(resp.Body).Decode(&geese[0]); err != nil {
		return "", err
	}
	if len(geese) < 1 {
		return "", fmt.Errorf("no geese in response from %s", uri)
	}
	a := geese[0]
	if a.Images.Small == "" {
		return "", fmt.Errorf("no image url in response from %s", uri)
	}
	// checking size, GitHub doesn't support big images
	toobig, err := scmTools.ImageTooBig(a.Images.Small)
	if err != nil {
		return "", fmt.Errorf("could not validate image size %s: %v", a.Images.Small, err)
	} else if toobig {
		return "", fmt.Errorf("long goose is too long: %s", a.Images.Small)
	}
	return a.Format()
}

func handleIssueComment(request plugins.PluginRequest, event *scm.IssueCommentHook) error {
	scmClient := request.ScmClient()
	setKey := func() {
		honk.setKey(request.Client(), request.RepoConfig().Namespace, request.PluginConfig().Goose.Key)
	}
	return handle(scmClient.Issues, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.Issue.Number, honk, setKey)
}

func handlePullRequestComment(request plugins.PluginRequest, event *scm.PullRequestCommentHook) error {
	scmClient := request.ScmClient()
	setKey := func() {
		honk.setKey(request.Client(), request.RepoConfig().Namespace, request.PluginConfig().Goose.Key)
	}
	return handle(scmClient.PullRequests, scmClient.Tools, request.Logger(), event.Repo, event.Action, event.Comment, event.PullRequest.Number, honk, setKey)
}

func handle(client scmClient, scmTools scmTools, logger logr.Logger, repo scm.Repository, action scm.Action, comment scm.Comment, number int, g gaggle, setKey func()) error {
	// Only consider new comments.
	if action != scm.ActionCreate {
		return nil
	}
	// Make sure they are requesting a goose
	mat := match.FindStringSubmatch(comment.Body)
	if mat == nil {
		return nil
	}

	// Now that we know this is a relevant event we can set the key.
	setKey()

	for i := 0; i < 3; i++ {
		resp, err := g.readGoose(scmTools)
		if err != nil {
			logger.Error(err, "Failed to get goose img")
			continue
		}
		return client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, resp))
	}

	msg := "Unable to find goose. Have you checked the garden?"
	if err := client.CreateComment(repo.FullName, number, plugins.FormatCommentResponse(scmTools, comment, msg)); err != nil {
		logger.Error(err, "Failed to leave comment")
	}

	return errors.New("could not find a valid goose image")
}
