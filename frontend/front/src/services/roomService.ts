import { BACKEND_URL } from "../constants";

const getRooms = async (token: string) => {
  const res = await fetch(`${BACKEND_URL}/api/rooms`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + token,
    },
  });

  return res;
};

const createRoom = async (token: string, name: string) => {
  const res = await fetch(`${BACKEND_URL}/api/rooms`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + token,
    },
    body: JSON.stringify({ name }),
  });

  return res;
};

export { getRooms, createRoom };
