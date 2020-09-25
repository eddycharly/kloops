package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	_ "github.com/eddycharly/kloops/pkg/chatbot/pluginimports"
	plugins "github.com/eddycharly/kloops/pkg/chatbot/plugins"
)

func main() {
	tmpl := `
# KLoops plugins

The list below lists the KLoops plugins.

Every plugin receive events from scm providers and react to those events in a specific way.

In addition, plugins can react to comments (on pull requests, issues or reviews). The list of supported commands is document [here](./COMMANDS.md).
{{ range $name, $plugin := . }}
- [{{ $name }}](#{{ $name }})
{{- end }}

{{ range $name, $plugin := . }}
## {{ $name }}

{{ $plugin.Description }}

This plugin reacts to the following events:
{{- range (events $plugin) }}
- {{ code . }}
{{ end }}

{{ if $plugin.ExcludedProviders -}}
This plugin does not support the following scm providers:
{{ range $plugin.ExcludedProviders }}
- {{ . }}
{{ end }}
{{- else -}}
This plugins supports all scm providers.
{{- end }}

{{ if $plugin.Commands -}}
This plugin has the following commands:
{{- range $plugin.Commands }}
{{- $command := . }}
{{- range (entries .Name) }}
- {{ cmd $command.Prefix . | code }}
{{- end }}
{{-  end }}
{{- else -}}
This plugins has no commands.
{{- end }}
{{ end }}
`

	t := template.New("doc").Funcs(template.FuncMap{
		"events": func(plugin plugins.Plugin) []string {
			ret := plugin.GetEvents()
			if ret == nil {
				ret = append(ret, "none")
			}
			return ret
		},
		"cmd": func(prefix, cmd string) string {
			ret := "/"
			if prefix != "" {
				ret += "[" + prefix + "]"
			}
			ret += cmd
			return ret
		},
		"code": func(code string) string {
			if code == "" {
				return ""
			}
			return fmt.Sprintf("`%s`", code)
		},
		"entries": func(name string) []string {
			return strings.Split(name, "|")
		},
	})
	if _, err := t.Parse(tmpl); err != nil {
		panic(err)
	}

	if err := t.Execute(os.Stdout, plugins.GetPlugins()); err != nil {
		panic(err)
	}

}
