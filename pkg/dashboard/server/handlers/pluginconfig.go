package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/dashboard/server/models"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginConfigHandler struct {
	namespace string
	client    client.Client
	logger    logr.Logger
}

func NewPluginConfigHandler(namespace string, client client.Client, logger logr.Logger) *PluginConfigHandler {
	logger = logger.WithName("PluginConfigHandler")
	return &PluginConfigHandler{
		namespace: namespace,
		client:    client,
		logger:    logger,
	}
}

func (h *PluginConfigHandler) List(w http.ResponseWriter, r *http.Request) {
	var list v1alpha1.PluginConfigList
	if err := h.client.List(r.Context(), &list, client.InNamespace(h.namespace)); err != nil {
		h.logger.Error(err, "failed to list plugin config")
	} else {
		if err := json.NewEncoder(w).Encode(models.FromPluginsList(list.Items)); err != nil {
			h.logger.Error(err, "failed to write response")
		}
	}
}

func (h *PluginConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := types.NamespacedName{
		Namespace: h.namespace,
		Name:      mux.Vars(r)["name"],
	}
	var item v1alpha1.PluginConfig
	if err := h.client.Get(r.Context(), key, &item); err != nil {
		h.logger.Error(err, "failed to get plugin config")
	} else {
		if err := json.NewEncoder(w).Encode(models.FromPlugin(item)); err != nil {
			h.logger.Error(err, "failed to write response")
		}
	}
}

// func (h *PluginConfigHandler) Create(w http.ResponseWriter, r *http.Request) {
// 	var item v1alpha1.PluginConfig
// 	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
// 		h.logger.Error(err, "failed to read body")
// 	}
// 	item.Namespace = h.namespace
// 	if err := h.client.Create(r.Context(), &item); err != nil {
// 		h.logger.Error(err, "failed to create repoconfigs")
// 	}
// 	fmt.Printf("%+v\n", item)
// 	if err := json.NewEncoder(w).Encode(models.ConvertFrom(item)); err != nil {
// 		h.logger.Error(err, "failed to write response")
// 	}
// }
