package chatbot

import (
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/git"
	utilsrepoconfig "github.com/eddycharly/kloops/pkg/utils/repoconfig"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/jenkins-x/go-scm/scm"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type handler struct {
	client    client.Client
	logger    logr.Logger
	gitClient git.Client
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	repo := vars["repo"]
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      repo,
		},
	}
	logger := h.logger.WithValues("namespace", namespace, "repo", repo)
	if r.Method != http.MethodPost {
		logger.WithValues("method", r.Method).Info("invalid http method so returning 200")
		return
	}
	var repoConfig v1alpha1.RepoConfig
	if err := h.client.Get(r.Context(), req.NamespacedName, &repoConfig); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Error(err, "resource not found in the cluster")
		} else {
			logger.Error(err, "error getting resource in the cluster")
		}
		return
	}
	webHook, err := utilsrepoconfig.ParseWebhook(r, &repoConfig)
	if err != nil {
		logger.Error(err, "failed to parse web hook")
		return
	}
	output, err := h.processWebHook(webHook)
	if err != nil {
		logger.Error(err, "failed to process web hook")
		return
	}
	repository := webHook.Repository()
	clone, err := utilsrepoconfig.GitClone(&repoConfig, h.gitClient, repository.FullName, repository.Clone)
	if err != nil {
		logger.Error(err, "failed to clone repository")
		return
	}
	fmt.Printf("%+v\n", clone)
	// // Demux events only to external plugins that require this event.
	// if external := util.ExternalPluginsForEvent(o.server.Plugins, string(webhook.Kind()), webhook.Repository().FullName); len(external) > 0 {
	// 	go util.CallExternalPluginsWithWebhook(l, external, webhook, util.HMACToken(), &o.server.wg)
	// }
	_, err = w.Write([]byte(output))
	if err != nil {
		logger.Error(err, "failed to write response")
		return
	}
	// fmt.Printf("%+v\n", webHook)
}

func (h *handler) processWebHook(webhook scm.Webhook) (string, error) {
	if _, ok := webhook.(*scm.PingHook); ok {
		return "pong", nil
	}
	// If we are in GitHub App mode and have a populated config, check if the repository for this webhook is one we actually
	// know about and error out if not.
	// if util.GetGitHubAppSecretDir() != "" && o.server.ConfigAgent != nil {
	// 	cfg := o.server.ConfigAgent.Config()
	// 	if cfg != nil {
	// 		if len(cfg.GetPostsubmits(repository)) == 0 && len(cfg.GetPresubmits(repository)) == 0 {
	// 			l.Infof("webhook from unconfigured repository %s, returning error", repository.Link)
	// 			return l, "", fmt.Errorf("repository not configured: %s", repository.Link)
	// 		}
	// 	}
	// }
	if _, ok := webhook.(*scm.PushHook); ok {
		// o.server.HandlePushEvent(l, pushHook)
		return "processed push hook", nil
	}
	if _, ok := webhook.(*scm.PullRequestHook); ok {
		// o.server.HandlePullRequestEvent(l, prHook)
		return "processed PR hook", nil
	}
	if _, ok := webhook.(*scm.BranchHook); ok {
		// o.server.HandleBranchEvent(l, branchHook)
		return "processed branch hook", nil
	}
	if _, ok := webhook.(*scm.IssueCommentHook); ok {
		// o.server.HandleIssueCommentEvent(l, *issueCommentHook)
		return "processed issue comment hook", nil
	}
	if _, ok := webhook.(*scm.PullRequestCommentHook); ok {
		// o.server.HandlePullRequestCommentEvent(l, *prCommentHook)
		return "processed PR comment hook", nil
	}
	if _, ok := webhook.(*scm.ReviewHook); ok {
		// o.server.HandleReviewEvent(l, *prReviewHook)
		return "processed PR review hook", nil
	}
	// l.Debugf("unknown kind %s webhook %#v", webhook.Kind(), webhook)
	return fmt.Sprintf("unknown hook %s", webhook.Kind()), nil
}

func Start(client client.Client) error {
	gitClient, err := git.NewClientWithDir("./tempo")
	if err != nil {
		return err
	}
	handler := handler{
		client:    client,
		logger:    ctrl.Log.WithName("chatbot"),
		gitClient: gitClient,
	}
	r := mux.NewRouter()
	r.Handle("/hook/{namespace}/{repo}", &handler)
	return http.ListenAndServe(":8090", r)
}
