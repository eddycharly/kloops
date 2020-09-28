export interface Secret {
  value?: string;
  valueFrom: ValueFrom;
}

export interface ValueFrom {
  secretKeyRef: SecretKeySelector;
}

export interface SecretKeySelector {
  name: string;
  key: string;
}

export interface PluginConfigSpec {
  owners: Owners;
  cat: Cat;
  goose: Goose;
  size: Size;
  welcome: Welcome;
}

export interface PluginConfig {
  name: string;
  namespace: string;
  creationTimestamp: string;
  spec: PluginConfigSpec;
}

export interface Cat {
  key: Secret;
}

export interface Goose {
  key: Secret;
}

export interface Size {
  s: number;
  m: number;
  l: number;
  xl: number;
  xxl: number;
}

export interface Welcome {
  messageTemplate?: string;
}

export interface Owners {
  // MDYAMLRepos []string `json:"mdyamlrepos,omitempty"`
  // // SkipCollaborators disables collaborator cross-checks and forces both
  // // the approve and lgtm plugins to use solely OWNERS files for access
  // // control in the provided repos.
  // SkipCollaborators []string `json:"skipCollaborators,omitempty"`
  // // LabelsExcludeList holds a list of labels that should not be present in any
  // // OWNERS file, preventing their automatic addition by the owners-label plugin.
  // // This check is performed by the verify-owners plugin.
  // LabelsExcludeList []string `json:"labelsExcludes,omitempty"`
}
