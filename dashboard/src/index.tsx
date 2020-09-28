import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import logger from 'redux-logger';
import { configureStore } from '@reduxjs/toolkit'
import './index.css';
import { App } from './App';
import { RootReducer } from 'reducers'

const store = configureStore({
  reducer: RootReducer,
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(logger),
});

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <App />
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);
