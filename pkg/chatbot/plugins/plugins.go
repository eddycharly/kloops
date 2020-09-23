/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugins

import (
	"fmt"

	"github.com/eddycharly/kloops/api/v1alpha1"
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
