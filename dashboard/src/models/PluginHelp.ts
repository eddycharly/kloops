export interface CommandArg {
  usage?: string;
  pattern: string;
  optional: string;
}

export interface Command {
  prefix?: string;
  names: Array<string>;
  arg: CommandArg;
  maxMatches?: number;
  description: string;
  whoCanUse?: string;
}

export interface PluginHelp {
  shortDescription: string;
  description: string;
  excludedProviders?: Array<string>;
  config?: { [name: string]: string };
  events?: Array<string>;
  commands?: Array<Command>;
}
