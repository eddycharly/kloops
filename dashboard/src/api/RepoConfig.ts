import { RepoConfig } from '../models/RepoConfig';
import { getApiEndpoint } from './api';

const path = '/api/repos';

export function listRepoConfigs(): Promise<Array<RepoConfig>> {
  const uri = `${getApiEndpoint()}${path}`;

  return fetch(uri)
    .then(response => {
      return response.json();
    });
}

export function getRepoConfig(name: string): Promise<RepoConfig> {
  const uri = `${getApiEndpoint()}${path}/${name}`;

  return fetch(uri)
    .then(response => {
      return response.json();
    });
}

export function createRepoConfig(config: any): Promise<any> {
  const uri = `${getApiEndpoint()}${path}`;

  return fetch(uri, {
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(config)
  }).then(response => {
    return response.json();
  });
}

export function createHook(name: any): Promise<any> {
  const uri = `${getApiEndpoint()}/api/hooks/${name}`;

  return fetch(uri, {
    method: 'POST',
  });
}
