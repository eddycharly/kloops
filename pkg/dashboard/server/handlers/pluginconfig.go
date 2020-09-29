package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/dashboard/server/models"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListPluginConfigs(namespace string, c client.Client, logger logr.Logger) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var list v1alpha1.PluginConfigList
		if err := c.List(r.Context(), &list, client.InNamespace(namespace)); err != nil {
			logger.Error(err, "failed to list plugin config")
		} else {
			if err := json.NewEncoder(w).Encode(models.FromPluginsList(list.Items)); err != nil {
				logger.Error(err, "failed to write response")
			}
		}
	}), nil
}

func GetPluginConfig(namespace string, c client.Client, logger logr.Logger) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := types.NamespacedName{
			Namespace: namespace,
			Name:      mux.Vars(r)["name"],
		}
		var item v1alpha1.PluginConfig
		if err := c.Get(r.Context(), key, &item); err != nil {
			logger.Error(err, "failed to get plugin config")
		} else {
			if err := json.NewEncoder(w).Encode(models.FromPlugin(item)); err != nil {
				logger.Error(err, "failed to write response")
			}
		}
	}), nil
}
