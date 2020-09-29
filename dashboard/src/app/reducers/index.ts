import { combineReducers } from '@reduxjs/toolkit';
import { Slice as pluginHelp } from './pluginHelp';

export const RootReducer = combineReducers({
    pluginHelp: pluginHelp.reducer,
});

export type RootState = ReturnType<typeof RootReducer>;
