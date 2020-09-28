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

const proxyRoute = "/proxy/{rest:.*}"

type Server interface {
	Start(addr string, port int) error
}

type server struct {
	namespace   string
	config      *rest.Config
	client      client.Client
	broadcaster *utils.Broadcaster
	logger      logr.Logger
}

func NewServer(namespace string, config *rest.Config, client client.Client, broadcaster *utils.Broadcaster, logger logr.Logger) Server {
	return &server{
		namespace:   namespace,
		config:      config,
		client:      client,
		broadcaster: broadcaster,
		logger:      logger.WithName("Server"),
	}
}

func (s *server) Start(addr string, port int) error {
	logger := s.logger.WithValues("addr", addr, "port", port)
	logger.Info("starting server ...")
	r := mux.NewRouter()
	// Proxy
	r.Handle(proxyRoute, handlers.NewProxyHandler(s.config, logger))
	// Api
	repoConfig := handlers.NewReponConfigHandler(s.namespace, s.client, s.logger)
	pluginConfig := handlers.NewPluginConfigHandler(s.namespace, s.client, s.logger)
	pluginHelp := handlers.NewPluginHelpHandler(s.logger)
	r.HandleFunc("/api/pluginhelp", pluginHelp.List).Methods("GET")
	r.HandleFunc("/api/plugins", pluginConfig.List).Methods("GET")
	r.HandleFunc("/api/plugins/{name}", pluginConfig.Get).Methods("GET")
	// r.HandleFunc("/api/plugins", pluginConfig.Create).Methods("POST")
	r.HandleFunc("/api/repos", repoConfig.List).Methods("GET")
	r.HandleFunc("/api/repos/{name}", repoConfig.Get).Methods("GET")
	r.HandleFunc("/api/repos", repoConfig.Create).Methods("POST")
	r.HandleFunc("/api/hooks/{name}", repoConfig.Hook).Methods("POST")
	// Websocket
	r.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		connection, err := utils.UpgradeToWebsocket(w, r)
		if err != nil {
			fmt.Println(err)
			// log.Error().Err(err).Msg("Could not upgrade to websocket connection")
			return
		}
		utils.WriteOnlyWebsocket(connection, s.broadcaster)
	})
	// Static content
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./dashboard/build"))))
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), r)
}
