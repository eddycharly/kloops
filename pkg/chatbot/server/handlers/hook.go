package handlers

import (
	"fmt"
	"net/http"

	"github.com/eddycharly/kloops/api/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/events"
	"github.com/eddycharly/kloops/pkg/git"
	"github.com/eddycharly/kloops/pkg/utils"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type hookHandler struct {
	namespace string
	client    client.Client
	logger    logr.Logger
	events    events.Events
}

func NewHookHandler(namespace string, client client.Client, logger logr.Logger) http.Handler {
	logger = logger.WithName("HookHandler")
	gitClient, err := git.NewClientWithDir("./tempo")
	if err != nil {
		logger.Error(err, "failed to create git client")
		fmt.Println(err)
	}
	return &hookHandler{
		namespace: namespace,
		client:    client,
		logger:    logger,
		events:    events.NewEvents(client, gitClient, logger),
	}
}

func (h *hookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	logger := h.logger.WithValues(utils.MapValues(vars)...).WithValues("method", r.Method)
	provider := vars["provider"]
	if service, err := factory.NewWebHookService(provider); err != nil {
		logger.Error(err, "error creating webhooks service")
	} else {
		var repoConfig v1alpha1.RepoConfig
		var nn types.NamespacedName
		if webhook, err := service.Parse(r, func(webhook scm.Webhook) (string, error) {
			nn = types.NamespacedName{
				Namespace: h.namespace,
				Name:      webhook.Repository().Name,
			}
			if err := h.client.Get(r.Context(), nn, &repoConfig); err != nil {
				if apierrors.IsNotFound(err) {
					logger.Error(err, "resource not found in the cluster")
				} else {
					logger.Error(err, "error getting resource in the cluster")
				}
				return "", err
			}
			// TODO gitea hard coded
			hmac, err := utils.GetSecret(h.client, repoConfig.Namespace, repoConfig.Spec.Gitea.HmacToken)
			if err != nil {
				logger.Error(err, "error getting hmac token")
				return "", err
			}
			return string(hmac), nil
		}); err != nil {
			logger.Error(err, "failed to parse web hook")
		} else {
			pluginConfigSpec := repoConfig.Spec.PluginConfig.Spec
			if pluginConfigSpec == nil {
				nn := types.NamespacedName{
					Namespace: repoConfig.Namespace,
					Name:      repoConfig.Spec.PluginConfig.Ref,
				}
				var pluginConfig v1alpha1.PluginConfig
				if err := h.client.Get(r.Context(), nn, &pluginConfig); err != nil {
					logger.Error(err, "error getting resource in the cluster")
				} else {
					pluginConfigSpec = &pluginConfig.Spec
				}
			}
			if pluginConfigSpec != nil {
				token, err := utils.GetSecret(h.client, repoConfig.Namespace, repoConfig.Spec.Gitea.Token)
				if err != nil {
					logger.Error(err, "error getting oauth token")
				} else {
					scmClient, err := factory.NewClient(provider, repoConfig.Spec.Gitea.ServerURL, string(token))
					if err != nil {
						logger.Error(err, "error creating scm client")
					} else {
						if output, err := h.events.ProcessWebHook(&repoConfig, pluginConfigSpec, scmClient, webhook); err != nil {
							logger.Error(err, "failed to process web hook")
						} else {
							if _, err = w.Write([]byte(output)); err != nil {
								logger.Error(err, "failed to write response")
							}
						}
					}
				}
			} else {
				logger.Info("no plugin config spec")
			}
		}
	}
}
