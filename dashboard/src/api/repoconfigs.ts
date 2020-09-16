import { getApiEndpoint } from './api';

const path = '/proxy/apis/config.kloops.io/v1alpha1/repoconfigs';

export function getRepoConfigs(): Promise<Array<any>> {
  const uri = `${getApiEndpoint()}${path}`;

  return fetch(uri)
    .then(response => {
      const contentType = response.headers.get('content-type');
      if (contentType) {
        if (contentType.includes('text/plain')) {
          return response.text();
        }
        if (contentType.includes('application/json')) {
          return response.json();
        }
      }
      throw "Unknow content type";
    })
    .then(response => response.items);
}
