package models

import (
	"github.com/eddycharly/kloops/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PluginConfig struct {
	Name              string                    `json:"name,omitempty"`
	Namespace         string                    `json:"namespace,omitempty"`
	CreationTimestamp metav1.Time               `json:"creationTimestamp,omitempty"`
	Spec              v1alpha1.PluginConfigSpec `json:"spec,omitempty"`
}

func FromPlugin(from v1alpha1.PluginConfig) PluginConfig {
	return PluginConfig{
		Name:              from.Name,
		Namespace:         from.Namespace,
		CreationTimestamp: from.CreationTimestamp,
		Spec:              from.Spec,
	}
}

func FromPluginsList(from []v1alpha1.PluginConfig) []PluginConfig {
	var list []PluginConfig
	for _, i := range from {
		list = append(list, FromPlugin(i))
	}
	return list
}
