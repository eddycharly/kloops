import { PluginHelp } from '../models';
import { getApiEndpoint } from './api';

const path = '/api/pluginhelp';

export function getPluginHelp(): Promise<{ [name: string]: PluginHelp }> {
  const uri = `${getApiEndpoint()}${path}`;

  return fetch(uri)
    .then(response => {
      return response.json();
    });
}
