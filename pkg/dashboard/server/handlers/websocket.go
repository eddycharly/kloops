package handlers

import (
	"net/http"

	"github.com/eddycharly/kloops/pkg/utils"
	"github.com/go-logr/logr"
)

func WebSocket(broadcaster *utils.Broadcaster, logger logr.Logger) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if connection, err := utils.UpgradeToWebsocket(w, r); err != nil {
			logger.Error(err, "Could not upgrade to websocket connection")
		} else {
			utils.WriteOnlyWebsocket(connection, broadcaster)
		}
	}), nil
}
