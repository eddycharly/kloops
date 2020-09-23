export interface Command {
    usage: string;
    featured: boolean;
    description: string;
    examples: Array<string>;
    whoCanUse: string;
}

export interface PluginHelp {
    shortDescription: string;
    description: string;
    excludedProviders?: Array<string>;
    config?: { [name: string]: string };
    events?: Array<string>;
    commands?: Array<Command>;
}

export interface Help {
    allRepos: Array<string>;
    repoPlugins: { [name: string]: string[] };
    repoExternalPlugins: { [name: string]: string[] };
    pluginHelp: { [name: string]: PluginHelp };
    externalPluginHelp: { [name: string]: PluginHelp };
}
