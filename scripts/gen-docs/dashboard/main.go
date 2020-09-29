package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/eddycharly/kloops/pkg/dashboard/server"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	tmpl := `
# KLoops Dashboard backend

## Routes

The routes below are the routes supported by the dashboard backend.

They are evaluated in order, the first match serves the request.

{{- range . }}

### {{ title . }}

{{ .Description }}

{{ if .Methods -}}
This route supports the following methods:
{{ range .Methods }}
- {{ . }}
{{ end }}
{{- else -}}
This route supports all methods.
{{- end }}

{{- end }}
`

	t := template.New("doc").
		Funcs(
			template.FuncMap{
				// 	"cmd": func(prefix, cmd string) string {
				// 		ret := "/"
				// 		if prefix != "" {
				// 			ret += "[" + prefix + "]"
				// 		}
				// 		ret += cmd
				// 		return ret
				// 	},
				// 	"arg": func(arg *plugins.CommandArg) string {
				// 		if arg == nil {
				// 			return ""
				// 		}
				// 		if arg.Optional {
				// 			return "Optional"
				// 		}
				// 		return "Mandatory"
				// 	},
				"code": func(code string) string {
					if code == "" {
						return ""
					}
					return fmt.Sprintf("`%s`", code)
				},
				"title": func(route server.Route) string {
					if route.PathPrefix != "" {
						return fmt.Sprintf("%s (prefix)", route.PathPrefix)
					}
					return route.Path
				},
				// 	"entries": func(name string) []string {
				// 		return strings.Split(name, "|")
				// 	},
				// 	"examples": func(prefix, cmd string) []string {
				// 		var ret []string
				// 		ret = append(ret, "/"+cmd)
				// 		if prefix != "" {
				// 			ret = append(ret, "/"+prefix+cmd)
				// 		}
				// 		return ret
				// 	},
			},
		)

	if _, err := t.Parse(tmpl); err != nil {
		panic(err)
	}
	s := server.NewServer("", nil, nil, nil, zap.New(zap.UseDevMode(true)))
	if err := t.Execute(os.Stdout, s.Routes); err != nil {
		panic(err)
	}
}
