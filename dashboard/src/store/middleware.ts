import ReconnectingWebSocket from "reconnecting-websocket";
import { Middleware } from 'redux';

export function createWebSocketMiddleware(socket: ReconnectingWebSocket): Middleware {
  return ({ dispatch }) => {
    socket.addEventListener('close', () => {
      dispatch({ type: 'WEBSOCKET_DISCONNECTED' });
    });

    socket.addEventListener('open', () => {
      dispatch({ type: 'WEBSOCKET_CONNECTED' });
    });

    socket.addEventListener('message', event => {
      if (event.type !== 'message') {
        return;
      }
      const message = JSON.parse(event.data);
      dispatch({ type: message.MessageType, payload: message.Payload });
    });
    return next => action => next(action);
  };
};
