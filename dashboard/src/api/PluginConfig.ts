import { PluginConfig } from '../models/PluginConfig';
import { getApiEndpoint } from './api';

const path = '/api/plugins';

export function listPluginConfigs(): Promise<Array<PluginConfig>> {
    const uri = `${getApiEndpoint()}${path}`;
  
    return fetch(uri)
      .then(response => {
        return response.json();
      });
  }
  
  export function getPluginConfig(name: string): Promise<PluginConfig> {
    const uri = `${getApiEndpoint()}${path}/${name}`;
  
    return fetch(uri)
      .then(response => {
        return response.json();
      });
  }
  