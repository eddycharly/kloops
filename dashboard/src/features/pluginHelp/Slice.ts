import { createAsyncThunk, createSlice, SerializedError } from '@reduxjs/toolkit'
import { PluginHelp } from 'models';
import { getPluginHelp } from 'api';

type LoadingState = { state: 'loading' };
type FinishedState = { state: 'finished'; data: { [name: string]: PluginHelp } };
type FailedState = { state: 'failed'; data: SerializedError };

type State = LoadingState | FinishedState | FailedState | null;

const initialState: State = { state: 'loading' };

export const FetchAll = createAsyncThunk(
  'PluginHelp/FetchAll',
  async () => {
    return await getPluginHelp();
  }
)

export const Slice = createSlice({
  name: 'PluginHelp',
  initialState: null as State,
  reducers: {},
  extraReducers: builder => {
    builder.addCase(FetchAll.pending, () => {
      return {
        state: 'loading'
      }
    })
    builder.addCase(FetchAll.fulfilled, (state, action) => {
      return {
        state: 'finished',
        data: action.payload
      }
    })
    builder.addCase(FetchAll.rejected, (state, action) => {
      return {
        state: 'failed',
        data: action.error,
      }
    })
  }
});
