package handlers

import (
	"net/http"
	"net/url"

	"github.com/eddycharly/kloops/pkg/proxy"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
)

type proxyHandler struct {
	logger logr.Logger
	client *http.Client
	host   string
}

func NewProxyHandler(config *rest.Config, logger logr.Logger) http.Handler {
	logger = logger.WithName("HookHandler")
	transport, err := rest.TransportFor(config)
	if err != nil {
		// log.Error().Err(err).Msg("Error building http transport")
	}
	return &proxyHandler{
		logger: logger,
		client: &http.Client{Transport: transport},
		host:   config.Host,
	}
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		return
	}
	uri := mux.Vars(r)["rest"] + "?" + parsedURL.RawQuery

	proxy.Proxy(r, w, h.host+"/"+uri, h.client)
}
