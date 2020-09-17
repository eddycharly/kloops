package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/dashboard/server/models"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RepoconfigHandler struct {
	logger logr.Logger
	client client.Client
}

func NewReponConfigHandler(client client.Client, logger logr.Logger) *RepoconfigHandler {
	logger = logger.WithName("RepoconfigHandler")
	return &RepoconfigHandler{
		logger: logger,
		client: client,
	}
}

func (h *RepoconfigHandler) List(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("List")
	var list v1alpha1.RepoConfigList
	if err := h.client.List(r.Context(), &list, client.InNamespace(mux.Vars(r)["namespace"])); err != nil {
		h.logger.Error(err, "failed to list repoconfigs")
	} else {
		if err := json.NewEncoder(w).Encode(models.ConvertList(list.Items)); err != nil {
			h.logger.Error(err, "failed to write response")
		}
	}
}
