import { configureStore, createSlice } from "@reduxjs/toolkit";
import type { TokenState, UserState } from "../types";

const initialUserState: UserState = {
  user: undefined,
};

const userSlice = createSlice({
  name: "user",
  initialState: initialUserState,
  reducers: {
    setUser: (state, action) => {
      state.user = action.payload;
    },
    logout: (state) => {
      state.user = undefined;
    },
  },
});

const initialTokenState: TokenState = {
  token: localStorage.getItem("JWT")
}

const tokenSlice = createSlice({
  name: "token",
  initialState: initialTokenState,
  reducers: {
    setToken: (state, action) => {
      localStorage.setItem("JWT", action.payload)
      state.token = action.payload;
    },
    clearToken: (state) => {
      state.token = null;
      localStorage.removeItem("JWT")
    },
  },
});

export const { setUser, logout } = userSlice.actions;
export const { setToken, clearToken } = tokenSlice.actions;

export const store = configureStore({
  reducer: {
    user: userSlice.reducer,
    token: tokenSlice.reducer,
  },
  
});

export type RootState = ReturnType<typeof store.getState>
