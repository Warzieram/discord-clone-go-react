import { BACKEND_URL } from "../constants";

const regsiter = async (email: string, password: string, username: string) => {
  const response = await fetch(`${BACKEND_URL}/api/register`, {
    method: "post",
    headers: {
      "Content-Type": "aplication/json",
    },
    body: JSON.stringify({
      email,
      password,
      username,
    }),
  });

  return response;
};

const login = async (email: string, password: string) => {
  const response = await fetch(`${BACKEND_URL}/api/login`, {
    method: "post",
    headers: {
      "Content-Type": "aplication/json",
    },
    body: JSON.stringify({
      email,
      password,
    }),
  });

  return response;
};

export { regsiter, login };
