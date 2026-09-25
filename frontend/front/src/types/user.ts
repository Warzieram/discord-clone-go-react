export type User = {
  id: number | undefined;
  email: string | undefined;
  created_at: string | undefined;
  username: string | undefined;
};

export type UserState = {
  user: User | undefined;
};

export type TokenState = {
  token: string | null;
};
