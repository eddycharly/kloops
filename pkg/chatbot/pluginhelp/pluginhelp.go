package pluginhelp

// CommandArg is a serializable representation of the command argument information for a single command.
type CommandArg struct {
	Usage    string `json:"usage,omitempty"`
	Pattern  string `json:"pattern"`
	Optional bool   `json:"optional"`
}

// Command is a serializable representation of the command information for a single command.
type Command struct {
	Prefix      string      `json:"prefix,omitempty"`
	Names       []string    `json:"names"`
	Arg         *CommandArg `json:"arg,omitempty"`
	MaxMatches  int         `json:"maxMatches,omitempty"`
	Description string      `json:"description"`
	WhoCanUse   string      `json:"whoCanUse"`
}

// PluginHelp is a serializable representation of the help information for a single plugin.
// This includes repo specific configuration for every repo that the plugin is enabled for.
type PluginHelp struct {
	ShortDescription  string            `json:"shortDescription,omitempty"`
	Description       string            `json:"description"`
	ExcludedProviders []string          `json:"excludedProviders,omitempty"`
	Config            map[string]string `json:"config,omitempty"`
	Events            []string          `json:"events,omitempty"`
	Commands          []Command         `json:"commands,omitempty"`
}

// AddCommand registers new help text for a bot command.
func (pluginHelp *PluginHelp) AddCommand(command Command) {
	pluginHelp.Commands = append(pluginHelp.Commands, command)
}
