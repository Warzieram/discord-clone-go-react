import type { User } from "./user";

/** Shape returned by both the login and register endpoints. */
export type AuthApiResponse = {
  token: string;
  user: User;
};
