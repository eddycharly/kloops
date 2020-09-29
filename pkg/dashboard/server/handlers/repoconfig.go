package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/dashboard/server/models"
	"github.com/eddycharly/kloops/pkg/utils"
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

func ListRepoConfigs(namespace string, c client.Client, logger logr.Logger) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var list v1alpha1.RepoConfigList
		if err := c.List(r.Context(), &list, client.InNamespace(namespace)); err != nil {
			logger.Error(err, "failed to list repoconfigs")
		} else {
			if err := json.NewEncoder(w).Encode(models.ConvertList(list.Items)); err != nil {
				logger.Error(err, "failed to write response")
			}
		}
	}), nil
}

func GetRepoConfig(namespace string, c client.Client, logger logr.Logger) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := types.NamespacedName{
			Namespace: namespace,
			Name:      mux.Vars(r)["name"],
		}
		var item v1alpha1.RepoConfig
		if err := c.Get(r.Context(), key, &item); err != nil {
			logger.Error(err, "failed to list repoconfigs")
		} else {
			if err := json.NewEncoder(w).Encode(models.ConvertFrom(item)); err != nil {
				logger.Error(err, "failed to write response")
			}
		}
	}), nil
}

func CreateRepoConfig(namespace string, c client.Client, logger logr.Logger) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var item v1alpha1.RepoConfig
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			logger.Error(err, "failed to read body")
		}
		item.Namespace = namespace
		if err := c.Create(r.Context(), &item); err != nil {
			logger.Error(err, "failed to create repoconfigs")
		}
		fmt.Printf("%+v\n", item)
		if err := json.NewEncoder(w).Encode(models.ConvertFrom(item)); err != nil {
			logger.Error(err, "failed to write response")
		}
	}), nil
}

func HookRepoConfig(namespace string, c client.Client, logger logr.Logger) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := types.NamespacedName{
			Namespace: namespace,
			Name:      mux.Vars(r)["name"],
		}
		var item v1alpha1.RepoConfig
		if err := c.Get(r.Context(), key, &item); err != nil {
			logger.Error(err, "failed to get repoconfig")
		} else {
			if scmClient, _, err := utilsrepoconfig.ScmClient(c, &item); err != nil {
				logger.Error(err, "failed to get scm client")
			} else {
				// TODO not always gitea...
				hmac, _ := utils.GetSecret(c, item.Namespace, item.Spec.Gitea.HmacToken)
				hook := scm.HookInput{
					Target:     "http://kloops-chatbot.tools.svc.cluster.local/hook/gitea/",
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
					logger.Error(err, "failed to create hook")
				}
			}
		}
	}), nil
}
