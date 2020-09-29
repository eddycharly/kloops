import { configureStore, Action } from '@reduxjs/toolkit';
import { ThunkAction } from 'redux-thunk';
import { RootReducer, RootState } from 'app/reducers';
import ReconnectingWebSocket from 'reconnecting-websocket';
import { createWebSocketMiddleware } from './middleware';
import { logger } from 'redux-logger';
import { getWebSocketEndpoint } from 'api';

export const WebSocket = new ReconnectingWebSocket(getWebSocketEndpoint());
export const Store = configureStore({
  reducer: RootReducer,
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(logger, createWebSocketMiddleware(WebSocket)),
});

export type AppDispatch = typeof Store.dispatch;
export type AppThunk = ThunkAction<void, RootState, unknown, Action<string>>;
