package server

import (
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/pkg/chatbot/server/handlers"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const helpRoute = "/help/{namespace}/{repo}"
const hookRoute = "/hook/{provider}"

type Server interface {
	Start(addr string, port int) error
}

type server struct {
	namespace string
	client    client.Client
	logger    logr.Logger
}

func NewServer(namespace string, client client.Client, logger logr.Logger) Server {
	return &server{
		namespace: namespace,
		client:    client,
		logger:    logger.WithName("Server"),
	}
}

func (s *server) Start(addr string, port int) error {
	logger := s.logger.WithValues("addr", addr, "port", port)
	logger.Info("starting server ...")
	r := mux.NewRouter()
	r.Handle(helpRoute, handlers.NewHelpHandler(s.client, logger))
	r.Handle(hookRoute, handlers.NewHookHandler(s.namespace, s.client, logger)).Methods("POST")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), r)
}
