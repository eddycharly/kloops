import {
  PluginConfigSpec,
  Secret
} from './PluginConfig';

export interface AutoMerge {
  batchSizeLimit: number;
  mergeType: string;
  labels: Array<string>;
  missingLabels: Array<string>;
  reviewApprovedRequired: boolean;
}

export interface GitHubRepo {
  owner: string;
  repo: string;
  server: string;
  hmacToken: Secret;
  token: Secret;
}

export interface GiteaRepo {
  owner: string;
  repo: string;
  server: string;
  hmacToken: Secret;
  token: Secret;
}

export interface RepoPluginConfig {
  ref?: string;
  spec: PluginConfigSpec;
  plugins: Array<string>;
}

export interface RepoConfigSpec {
  botName: string;
  gitHub: GitHubRepo;
  gitea: GiteaRepo;
  autoMerge: AutoMerge;
  pluginConfig: RepoPluginConfig;
}

export interface RepoConfig {
  name: string;
  namespace: string
  spec: RepoConfigSpec;
}
