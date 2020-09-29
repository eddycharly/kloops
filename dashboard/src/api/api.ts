export function getApiEndpoint(): string {
    return window.location.origin;
}

export function getWebSocketEndpoint(): string {
  const apiEndpoint = getApiEndpoint();
  return `${apiEndpoint.replace(/^http/, 'ws')}/resources`;
}
