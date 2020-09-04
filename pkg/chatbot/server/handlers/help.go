package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	utils "github.com/eddycharly/kloops/pkg/utils/logr"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type helpHandler struct {
	client client.Client
	logger logr.Logger
}

func NewHelpHandler(client client.Client, logger logr.Logger) http.Handler {
	return &helpHandler{
		client: client,
		logger: logger.WithName("HelpHandler"),
	}
}

func (h *helpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	logger := h.logger.WithValues(utils.MapValues(vars)...).WithValues("method", r.Method)
	if r.Method != http.MethodGet {
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
				help := make(map[string]*pluginhelp.PluginHelp)
				for k, v := range plugins.HelpProviders() {
					if h, err := v(pluginConfigSpec); err != nil {
						logger.WithValues("plugin", k).Error(err, "failed to retrieve plugin help")
					} else {
						help[k] = h
					}
				}
				if b, err := json.Marshal(help); err != nil {
					logger.Error(err, "failed to marshaling plugin help")
				} else {
					fmt.Fprint(w, string(b))
				}
			} else {
				logger.Info("no plugin config spec")
			}
		}
	}
}
