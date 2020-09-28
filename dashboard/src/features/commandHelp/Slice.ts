import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { PluginHelp } from 'models';
import { getPluginHelp } from 'api';

type LoadingState = { state: 'loading' };
type FinishedState = { state: 'finished'; data: { [name: string]: PluginHelp } };
type FailedState = { state: 'failed'; data: Error };

type State = LoadingState | FinishedState | FailedState;

const initialState: State = { state: 'loading' };

export const FetchAll = createAsyncThunk(
  'CommandHelp/FetchAll',
  async () => {
    return await getPluginHelp();
  }
)

export const Slice = createSlice({
  name: 'CommandHelp',
  initialState: initialState as State,
  reducers: {},
  extraReducers: builder => {
    builder.addCase(FetchAll.fulfilled, (state, action) => {
      return {
        state: 'finished',
        data: action.payload
      }
    })
  }
});
