import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { App } from 'app/App';
import { SnackbarProvider } from 'notistack';
import { Store, WebSocket } from 'app/store';
import './index.css';

function closeSocket() {
  WebSocket.close();
}

ReactDOM.render(
  <React.StrictMode>
    <SnackbarProvider maxSnack={3}>
      <Provider store={Store}>
        <App onUnload={closeSocket} />
      </Provider>
    </SnackbarProvider>
  </React.StrictMode>,
  document.getElementById('root')
);
