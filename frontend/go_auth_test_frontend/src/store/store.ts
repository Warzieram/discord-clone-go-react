import { configureStore, createSlice } from "@reduxjs/toolkit";

export type User = {
  id: number | undefined;
  email: string | undefined;
  created_at: string | undefined;
};

type UserState = {
  user: User | undefined;
};

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

type TokenState = {
  token: string | null
}
const initialTokenState: TokenState = {
  token: localStorage.getItem("JWT")
}

const tokenSlice = createSlice({
  name: "token",
  initialState: initialTokenState,
  reducers: {
    setToken: (state, action) => {
      localStorage.setItem("JWT", action.payload)
      state = action.payload;
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
