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
# KLoops commands

The table below lists the commands KLoops understands.

In addition, you can prefix commands with {{ code "kl-" }}. For example, {{ code "/meow" }} and {{ code "/kl-meow" }} are equivalent.

This is because Gitlab hijacks some slash commands for its own [quick actions](https://docs.gitlab.com/ee/user/project/quick_actions.html), and we never get notified.
In practice, we don’t need the {{ code "kl-" }} prefix for everything, just commands that also are quick actions,
but we opted to play it safe and add {{ code "kl-" }} prefixes for every command just in case Gitlab eventually adds conflicting quick actions.

There is a [PR open](https://gitlab.com/gitlab-org/gitlab/-/issues/215934) for them to send webhook events for quick actions,
at which point we wouldn't need to worry about it any more, but for now, we do need the {{ code "/kl-(command)" }}.


| Command | Argument | Description | Examples | Plugin |
|---------|----------|-------------|----------|--------|
{{- range $name, $plugin := . }}
{{- range $plugin.Commands }}
{{- $command := . }}
{{- range (entries .Name) }}
{{- $examples := examples $command.Prefix . }}
| {{ cmd $command.Prefix . | code }} | {{ arg $command.Arg }} | {{ $command.Description }} | <ul>{{- range $examples }}<li>{{ code . }}</li>{{ end }}</ul> | [{{ $name }}](./PLUGINS.md#{{ $name }}) |
{{- end }}
{{- end }}
{{- end }}`

	t := template.New("doc").Funcs(template.FuncMap{
		"cmd": func(prefix, cmd string) string {
			ret := "/"
			if prefix != "" {
				ret += "[" + prefix + "]"
			}
			ret += cmd
			return ret
		},
		"arg": func(arg *plugins.CommandArg) string {
			if arg == nil {
				return ""
			}
			if arg.Optional {
				return "Optional"
			}
			return "Mandatory"
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
		"examples": func(prefix, cmd string) []string {
			var ret []string
			ret = append(ret, "/"+cmd)
			if prefix != "" {
				ret = append(ret, "/"+prefix+cmd)
			}
			return ret
		},
	})
	if _, err := t.Parse(tmpl); err != nil {
		panic(err)
	}

	if err := t.Execute(os.Stdout, plugins.GetPlugins()); err != nil {
		panic(err)
	}

}
