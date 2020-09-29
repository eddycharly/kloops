package server

import (
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/pkg/dashboard/server/handlers"
	"github.com/eddycharly/kloops/pkg/utils"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Route struct {
	Description    string
	Path           string
	PathPrefix     string
	Methods        []string
	HandlerFactory func() (http.Handler, error)
}

type Server interface {
	Start(addr string, port int) error
}

type DashboardServer struct {
	Routes []Route
	logger logr.Logger
}

func NewServer(namespace string, config *rest.Config, client client.Client, broadcaster *utils.Broadcaster, logger logr.Logger) DashboardServer {
	logger = logger.WithName("Server")
	// proxy := handlers.NewProxyHandler(config, logger)
	// static := http.StripPrefix("/", http.FileServer(http.Dir("./dashboard/build")))
	// repoConfig := handlers.NewReponConfigHandler(namespace, client, logger)
	// pluginConfig := handlers.NewPluginConfigHandler(namespace, client, logger)
	// pluginHelp := handlers.NewPluginHelpHandler(logger)
	var routes = []Route{
		{
			Description: "Proxy to the kubernetes api server",
			Path:        "/proxy/{rest:.*}",
			HandlerFactory: func() (http.Handler, error) {
				return handlers.Proxy(config, logger)
			},
		},
		{
			Description: "Get plugins help",
			Path:        "/api/pluginhelp",
			Methods:     []string{"GET"},
			HandlerFactory: func() (http.Handler, error) {
				return handlers.PluginHelp(logger)
			},
		},
		{
			Description: "List plugin configs",
			Path:        "/api/plugins",
			Methods:     []string{"GET"},
			HandlerFactory: func() (http.Handler, error) {
				return handlers.ListPluginConfigs(namespace, client, logger)
			},
		},
		{
			Description: "Get plugin config",
			Path:        "/api/plugins/{name}",
			Methods:     []string{"GET"},
			HandlerFactory: func() (http.Handler, error) {
				return handlers.GetPluginConfig(namespace, client, logger)
			},
		},
		{
			Description: "List repo configs",
			Path:        "/api/repos",
			Methods:     []string{"GET"},
			HandlerFactory: func() (http.Handler, error) {
				return handlers.ListRepoConfigs(namespace, client, logger)
			},
		},
		{
			Description: "Get repo config",
			Path:        "/api/repos/{name}",
			Methods:     []string{"GET"},
			HandlerFactory: func() (http.Handler, error) {
				return handlers.GetRepoConfig(namespace, client, logger)
			},
		},
		{
			Description: "Create repo config",
			Path:        "/api/repos/{name}",
			Methods:     []string{"POST"},
			HandlerFactory: func() (http.Handler, error) {
				return handlers.CreateRepoConfig(namespace, client, logger)
			},
		},
		{
			Description: "Setup repo config hooks",
			Path:        "/api/hooks/{name}",
			Methods:     []string{"POST"},
			HandlerFactory: func() (http.Handler, error) {
				return handlers.HookRepoConfig(namespace, client, logger)
			},
		},
		{
			Description: "Connect to websocket",
			Path:        "/socket",
			Methods:     []string{"GET"},
			HandlerFactory: func() (http.Handler, error) {
				return handlers.WebSocket(broadcaster, logger)
			},
		},
		{
			Description: "Static content",
			PathPrefix:  "/",
			Methods:     []string{"GET"},
			HandlerFactory: func() (http.Handler, error) {
				return http.StripPrefix("/", http.FileServer(http.Dir("./dashboard/build"))), nil
			},
		},
	}

	return DashboardServer{
		Routes: routes,
		logger: logger,
	}
}

func (s *DashboardServer) Start(addr string, port int) error {
	logger := s.logger.WithValues("addr", addr, "port", port)
	logger.Info("setting up server routes ...")
	router := mux.NewRouter()
	for _, route := range s.Routes {
		logger.WithValues("path", route.Path, "pathPrefix", route.PathPrefix, "description", route.Description).Info("setting up server routes ...")
		var r *mux.Route
		handler, err := route.HandlerFactory()
		if err != nil {
			logger.Error(err, "Error building http transport")
			return err
		}
		if route.PathPrefix != "" {
			r = router.PathPrefix(route.PathPrefix).Handler(handler)
		} else {
			r = router.Handle(route.Path, handler)
		}
		if route.Methods != nil {
			r.Methods(route.Methods...)
		}
	}
	logger.Info("starting server ...")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), router)
}
