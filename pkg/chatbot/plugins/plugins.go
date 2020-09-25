package plugins

import (
	"fmt"

	"github.com/eddycharly/kloops/apis/config/v1alpha1"
	"github.com/eddycharly/kloops/pkg/chatbot/pluginhelp"
)

var (
	plugins = map[string]Plugin{}
)

// RegisterPlugin registers a plugin.
func RegisterPlugin(name string, plugin Plugin) {
	plugins[name] = plugin
}

// HelpProvider defines the function type that construct a pluginhelp.PluginHelp for enabled
// plugins. It takes into account the plugins configuration and enabled repositories.
type HelpProvider func(*v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error)

// HelpProviders returns the map of registered plugins with their associated HelpProvider.
func HelpProviders() map[string]HelpProvider {
	pluginHelp := make(map[string]HelpProvider)
	for k, v := range plugins {
		h := v
		pluginHelp[k] = func(config *v1alpha1.PluginConfigSpec) (*pluginhelp.PluginHelp, error) {
			return h.GetHelp(config)
		}
	}
	return pluginHelp
}

// GetPlugin returns the registered plugin for a given name.
func GetPlugin(name string) (*Plugin, error) {
	if plugin, ok := plugins[name]; ok {
		return &plugin, nil
	}
	return nil, fmt.Errorf("plugin %s not found", name)
}

// GetPlugins returns the plugins map
func GetPlugins() map[string]Plugin {
	return plugins
}
