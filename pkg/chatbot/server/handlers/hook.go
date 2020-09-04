package handlers

import (
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/events"
	"github.com/eddycharly/kloops/pkg/git"
	utils "github.com/eddycharly/kloops/pkg/utils/logr"
	utilsrepoconfig "github.com/eddycharly/kloops/pkg/utils/repoconfig"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type hookHandler struct {
	client client.Client
	logger logr.Logger
	events events.Events
}

func NewHookHandler(client client.Client, logger logr.Logger) http.Handler {
	logger = logger.WithName("HookHandler")
	gitClient, err := git.NewClientWithDir("./tempo")
	if err != nil {
		logger.Error(err, "failed to create git client")
		fmt.Println(err)
	}
	return &hookHandler{
		client: client,
		logger: logger,
		events: events.NewEvents(client, gitClient, logger),
	}
}

func (h *hookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	logger := h.logger.WithValues(utils.MapValues(vars)...).WithValues("method", r.Method)
	if r.Method != http.MethodPost {
		logger.Info("invalid http method so returning 200")
	} else {
		nn := types.NamespacedName{
			Namespace: vars["namespace"],
			Name:      vars["repo"],
		}
		var repoConfig v1alpha1.RepoConfig
		if err := h.client.Get(r.Context(), nn, &repoConfig); err != nil {
			if apierrors.IsNotFound(err) {
				logger.Error(err, "resource not found in the cluster")
			} else {
				logger.Error(err, "error getting resource in the cluster")
			}
		} else {
			pluginConfigSpec := repoConfig.Spec.PluginConfig.Spec
			if pluginConfigSpec == nil {
				nn := types.NamespacedName{
					Namespace: repoConfig.Namespace,
					Name:      repoConfig.Spec.PluginConfig.Ref,
				}
				var pluginConfig v1alpha1.PluginConfig
				if err := h.client.Get(r.Context(), nn, &pluginConfig); err != nil {
					logger.Error(err, "error getting resource in the cluster")
				} else {
					pluginConfigSpec = &pluginConfig.Spec
				}
			}
			if pluginConfigSpec != nil {
				if event, scmClient, err := utilsrepoconfig.ParseWebhook(h.client, r, &repoConfig); err != nil {
					logger.Error(err, "failed to parse web hook")
				} else {
					if output, err := h.events.ProcessWebHook(&repoConfig, pluginConfigSpec, scmClient, event); err != nil {
						logger.Error(err, "failed to process web hook")
					} else {
						if _, err = w.Write([]byte(output)); err != nil {
							logger.Error(err, "failed to write response")
						}
					}
				}
			} else {
				logger.Info("no plugin config spec")
			}
		}
	}
}
