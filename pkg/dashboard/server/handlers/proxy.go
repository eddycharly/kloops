package handlers

import (
	"net/http"
	"net/url"

	"github.com/eddycharly/kloops/pkg/proxy"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
)

func Proxy(config *rest.Config, logger logr.Logger) (http.Handler, error) {
	transport, err := rest.TransportFor(config)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Transport: transport}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsedURL, err := url.Parse(r.URL.String())
		if err != nil {
			return
		}
		uri := mux.Vars(r)["rest"] + "?" + parsedURL.RawQuery

		proxy.Proxy(r, w, config.Host+"/"+uri, client)
	}), nil
}
