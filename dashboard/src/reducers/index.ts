import { combineReducers } from 'redux';
import { Slice as commandHelp } from 'features/commandHelp/Slice';
import { Slice as pluginHelp } from 'features/pluginHelp/Slice';

export const RootReducer = combineReducers({
    commandHelp: commandHelp.reducer,
    pluginHelp: pluginHelp.reducer,
});

export type RootState = ReturnType<typeof RootReducer>;
