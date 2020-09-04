package cat

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
	match          = regexp.MustCompile(`(?mi)^/(?:lh-)?meow(vie)?(?: (.+))?\s*$`)
	grumpyKeywords = regexp.MustCompile(`(?mi)^(no|grumpy)\s*$`)
	meow           = &realClowder{
		url: "https://api.thecatapi.com/v1/images/search?format=json&results_per_page=1",
	}
)

const (
	pluginName = "cat"
	grumpyURL  = "https://upload.wikimedia.org/wikipedia/commons/e/ee/Grumpy_Cat_by_Gage_Skidmore.jpg"
)

func init() {
	plugins.RegisterHelpProvider(pluginName, helpProvider)
	plugins.RegisterIssueCommentHandler(pluginName, handleIssueComment)
	plugins.RegisterPullRequestCommentHandler(pluginName, handlePullRequestComment)
}

func helpProvider(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
	pluginHelp := &pluginhelp.PluginHelp{
		Description: "The cat plugin adds a cat image to an issue or PR in response to the `/meow` command.",
		Config: map[string]string{
			"": "The cat plugin uses an api key for thecatapi.com stored in the plugin config.",
		},
	}
	pluginHelp.AddCommand(pluginhelp.Command{
		Usage:       "/meow(vie) [CATegory]",
		Description: "Add a cat image to the issue or PR",
		Featured:    false,
		WhoCanUse:   "Anyone",
		Examples:    []string{"/meow", "/meow caturday", "/meowvie clothes", "/lh-meow"},
	})
	return pluginHelp, nil
}

type scmClient interface {
	CreateComment(owner, repo string, number int, pr bool, comment string) error
	QuoteAuthorForComment(string) string
}

type clowder interface {
	readCat(string, bool) (string, error)
}

type realClowder struct {
	url    string
	lock   sync.RWMutex
	update time.Time
	key    string
}

func (c *realClowder) setKey(client client.Client, namespace string, secret v1alpha1.Secret) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !time.Now().After(c.update) {
		return
	}
	c.update = time.Now().Add(1 * time.Minute)
	key, err := utils.GetSecret(client, namespace, secret)
	if err == nil {
		c.key = strings.TrimSpace(string(key))
		return
	}
	// log.WithValues("keyPath", keyPath).Error(err, "failed to read key")
	c.key = ""
}

type catResult struct {
	Image string `json:"url"`
}

func (cr catResult) Format() (string, error) {
	if cr.Image == "" {
		return "", errors.New("empty image url")
	}
	img, err := url.Parse(cr.Image)
	if err != nil {
		return "", fmt.Errorf("invalid image url %s: %v", cr.Image, err)
	}

	return fmt.Sprintf("![cat image](%s)", img), nil
}

func (c *realClowder) URL(category string, movieCat bool) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	uri := string(c.url)
	if category != "" {
		uri += "&category=" + url.QueryEscape(category)
	}
	if c.key != "" {
		uri += "&api_key=" + url.QueryEscape(c.key)
	}
	if movieCat {
		uri += "&mime_types=gif"
	}
	return uri
}

func (c *realClowder) readCat(category string, movieCat bool) (string, error) {
	cats := make([]catResult, 0)
	uri := c.URL(category, movieCat)
	if grumpyKeywords.MatchString(category) {
		cats = append(cats, catResult{grumpyURL})
	} else {
		resp, err := http.Get(uri) // #nosec
		if err != nil {
			return "", fmt.Errorf("could not read cat from %s: %v", uri, err)
		}
		defer resp.Body.Close()
		if sc := resp.StatusCode; sc > 299 || sc < 200 {
			return "", fmt.Errorf("failing %d response from %s", sc, uri)
		}
		if err = json.NewDecoder(resp.Body).Decode(&cats); err != nil {
			return "", err
		}
		if len(cats) < 1 {
			return "", fmt.Errorf("no cats in response from %s", uri)
		}
	}
	a := cats[0]
	if a.Image == "" {
		return "", fmt.Errorf("no image url in response from %s", uri)
	}
	// checking size, GitHub doesn't support big images
	toobig, err := utils.ImageTooBig(a.Image)
	if err != nil {
		return "", fmt.Errorf("could not validate image size %s: %v", a.Image, err)
	} else if toobig {
		return "", fmt.Errorf("longcat is too long: %s", a.Image)
	}
	return a.Format()
}

func handleIssueComment(request plugins.PluginRequest, event *scm.IssueCommentHook) error {
	setKey := func() { meow.setKey(request.Client(), request.Namespace(), request.PluginConfig().Cat.Key) }
	return handle(request.ScmClient(), request.Logger(), event.Repo, event.Action, event.Comment, event.Issue.Number, false, meow, setKey)
}

func handlePullRequestComment(request plugins.PluginRequest, event *scm.PullRequestCommentHook) error {
	setKey := func() { meow.setKey(request.Client(), request.Namespace(), request.PluginConfig().Cat.Key) }
	return handle(request.ScmClient(), request.Logger(), event.Repo, event.Action, event.Comment, event.PullRequest.Number, true, meow, setKey)
}

func handle(client scmClient, logger logr.Logger, repo scm.Repository, action scm.Action, comment scm.Comment, number int, pr bool, c clowder, setKey func()) error {
	// Only consider new comments.
	if action != scm.ActionCreate {
		return nil
	}
	// Make sure they are requesting a cat
	mat := match.FindStringSubmatch(comment.Body)
	if mat == nil {
		return nil
	}

	category, movieCat, err := parseMatch(mat)
	if err != nil {
		return err
	}

	// Now that we know this is a relevant event we can set the key.
	setKey()

	for i := 0; i < 3; i++ {
		resp, err := c.readCat(category, movieCat)
		if err != nil {
			logger.Error(err, "Failed to get cat img")
			continue
		}
		logger.Info(resp)
		err = client.CreateComment(
			repo.Namespace,
			repo.Name,
			number,
			pr,
			plugins.FormatResponseRaw(comment.Body, comment.Link, client.QuoteAuthorForComment(comment.Author.Login), resp),
		)
		if err != nil {
			logger.Error(err, "Failed to create comment")
		}
		return err
	}
	msg := "https://thecatapi.com appears to be down"
	if category != "" {
		msg = "Bad category. Please see https://api.thecatapi.com/api/categories/list"
	}
	err = client.CreateComment(
		repo.Namespace,
		repo.Name,
		number,
		pr,
		plugins.FormatResponseRaw(comment.Body, comment.Link, client.QuoteAuthorForComment(comment.Author.Login), msg),
	)
	if err != nil {
		logger.Error(err, "Failed to leave comment")
	}
	return errors.New("could not find a valid cat image")
}

func parseMatch(mat []string) (string, bool, error) {
	if len(mat) != 3 {
		err := fmt.Errorf("expected 3 capture groups in regexp match, but got %d", len(mat))
		return "", false, err
	}
	category := strings.TrimSpace(mat[2])
	movieCat := len(mat[1]) > 0 // "vie" suffix is present.
	return category, movieCat, nil
}
