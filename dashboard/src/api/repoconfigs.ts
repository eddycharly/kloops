import { getApiEndpoint } from './api';

const path = '/api/repos';

export function listRepoConfigs(): Promise<Array<any>> {
  const uri = `${getApiEndpoint()}${path}`;

  return fetch(uri)
    .then(response => {
      return response.json();
    })
}
