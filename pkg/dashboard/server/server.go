package server

import (
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/pkg/dashboard/server/handlers"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const proxyRoute = "/proxy/{rest:.*}"

type Server interface {
	Start(addr string, port int) error
}

type server struct {
	config *rest.Config
	client client.Client
	logger logr.Logger
}

func NewServer(config *rest.Config, client client.Client, logger logr.Logger) Server {
	return &server{
		config: config,
		client: client,
		logger: logger.WithName("Server"),
	}
}

func (s *server) Start(addr string, port int) error {
	logger := s.logger.WithValues("addr", addr, "port", port)
	logger.Info("starting server ...")
	r := mux.NewRouter()
	// Proxy
	r.Handle(proxyRoute, handlers.NewProxyHandler(s.config, logger))
	// Api
	repoConfig := handlers.NewReponConfigHandler(s.client, s.logger)
	r.HandleFunc("/api/repos", repoConfig.List).Methods("GET")
	// r.HandleFunc("/repos/{namespace}/{name}", repoConfig.List).Methods("GET")
	// Static content
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./dashboard/build"))))
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), r)
}
