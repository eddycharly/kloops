package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/dashboard/server/models"
	utilsrepoconfig "github.com/eddycharly/kloops/pkg/utils/repoconfig"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/jenkins-x/go-scm/scm"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RepoconfigHandler struct {
	namespace string
	client    client.Client
	logger    logr.Logger
}

func NewReponConfigHandler(namespace string, client client.Client, logger logr.Logger) *RepoconfigHandler {
	logger = logger.WithName("RepoconfigHandler")
	return &RepoconfigHandler{
		namespace: namespace,
		client:    client,
		logger:    logger,
	}
}

func (h *RepoconfigHandler) List(w http.ResponseWriter, r *http.Request) {
	var list v1alpha1.RepoConfigList
	if err := h.client.List(r.Context(), &list, client.InNamespace(h.namespace)); err != nil {
		h.logger.Error(err, "failed to list repoconfigs")
	} else {
		if err := json.NewEncoder(w).Encode(models.ConvertList(list.Items)); err != nil {
			h.logger.Error(err, "failed to write response")
		}
	}
}

func (h *RepoconfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := types.NamespacedName{
		Namespace: h.namespace,
		Name:      mux.Vars(r)["name"],
	}
	var item v1alpha1.RepoConfig
	if err := h.client.Get(r.Context(), key, &item); err != nil {
		h.logger.Error(err, "failed to list repoconfigs")
	} else {
		if err := json.NewEncoder(w).Encode(models.ConvertFrom(item)); err != nil {
			h.logger.Error(err, "failed to write response")
		}
	}
}

func (h *RepoconfigHandler) Create(w http.ResponseWriter, r *http.Request) {
	var item v1alpha1.RepoConfig
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		h.logger.Error(err, "failed to read body")
	}
	item.Namespace = h.namespace
	if err := h.client.Create(r.Context(), &item); err != nil {
		h.logger.Error(err, "failed to create repoconfigs")
	}
	fmt.Printf("%+v\n", item)
	if err := json.NewEncoder(w).Encode(models.ConvertFrom(item)); err != nil {
		h.logger.Error(err, "failed to write response")
	}
}

func (h *RepoconfigHandler) Hook(w http.ResponseWriter, r *http.Request) {
	key := types.NamespacedName{
		Namespace: h.namespace,
		Name:      mux.Vars(r)["name"],
	}
	var item v1alpha1.RepoConfig
	if err := h.client.Get(r.Context(), key, &item); err != nil {
		h.logger.Error(err, "failed to list repoconfigs")
	} else {
		if scmClient, _, err := utilsrepoconfig.ScmClient(h.client, &item); err != nil {
			h.logger.Error(err, "failed to get scm client")
		} else {
			// TODO not always gitea...
			hmac, _ := utilsrepoconfig.GetSecret(h.client, item.Namespace, item.Spec.Gitea.HmacToken)
			hook := scm.HookInput{
				Target:     fmt.Sprintf("http://kloops-chatbot.tools.svc.cluster.local/hook/%s/%s", key.Namespace, key.Name),
				Name:       key.Name,
				Secret:     string(hmac),
				SkipVerify: true,
				Events: scm.HookEvents{
					Branch:             true,
					Issue:              true,
					IssueComment:       true,
					PullRequest:        true,
					PullRequestComment: true,
					Push:               true,
					ReviewComment:      true,
					Tag:                true,
				},
			}
			if _, _, err := scmClient.Repositories.CreateHook(r.Context(), fmt.Sprintf("%s/%s", item.Spec.Gitea.Owner, item.Spec.Gitea.Repo), &hook); err != nil {
				h.logger.Error(err, "failed to create hook")
			}
		}
	}
}
