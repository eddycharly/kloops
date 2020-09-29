package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
	"github.com/eddycharly/kloops/pkg/chatbot/plugins"
	"github.com/go-logr/logr"
)

func PluginHelp(logger logr.Logger) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		help := make(map[string]*pluginhelp.PluginHelp)
		var pluginConfig v1alpha1.PluginConfig
		for k, v := range plugins.HelpProviders() {
			if h, err := v(&pluginConfig.Spec); err != nil {
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
	}), nil
}
