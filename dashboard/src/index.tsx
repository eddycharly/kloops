import { React, ReactDOM } from 'utils/imports/react';
import { configureStore, logger, Provider, ReconnectingWebSocket, RootReducer } from 'utils/imports/redux';
import { getWebSocketEndpoint } from 'api';
import { App } from './App';
import { createWebSocketMiddleware } from 'store/middleware';
import './index.css';

const webSocket = new ReconnectingWebSocket(getWebSocketEndpoint());
function closeSocket() {
  webSocket.close();
}

const store = configureStore({
  reducer: RootReducer,
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(logger, createWebSocketMiddleware(webSocket)),
});

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <App onUnload={closeSocket} />
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);
