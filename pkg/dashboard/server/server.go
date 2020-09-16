package server

import (
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/pkg/dashboard/server/handlers"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
)

const proxyRoute = "/proxy/{rest:.*}"

type Server interface {
	Start(addr string, port int) error
}

type server struct {
	config *rest.Config
	logger logr.Logger
}

func NewServer(config *rest.Config, logger logr.Logger) Server {
	return &server{
		config: config,
		logger: logger.WithName("Server"),
	}
}

func (s *server) Start(addr string, port int) error {
	logger := s.logger.WithValues("addr", addr, "port", port)
	logger.Info("starting server ...")
	r := mux.NewRouter()
	r.Handle(proxyRoute, handlers.NewProxyHandler(s.config, logger))
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./dashboard/build"))))
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), r)
}
