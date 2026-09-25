import { BACKEND_URL, WS_BACKEND_URL } from "../constants";

const getLastMessages = async (
  token: string,
  roomId: number,
  limit: number,
  offset: number,
) => {
  const res = await fetch(
    `${BACKEND_URL}/api/messages?roomID=${roomId}&limit=${limit}&offset=${offset}`,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + token,
      },
    },
  );

  return res;
};

const createMessageWebSocket = (token: string, roomId: number) => {
  return new WebSocket(`${WS_BACKEND_URL}/api/message?roomID=${roomId}`, [
    `auth.${token}`,
  ]);
};

const sendMessageWS = (ws: WebSocket, input: string) => {
  const request = {
    command_type: "SEND",
    data: input,
  };
  ws.send(JSON.stringify(request));
};

const sendDeleteRequest = (id: number, ws: WebSocket | null) => {
  if (!ws) return;
  const request = {
    command_type: "REMOVE",
    data: id,
  };
  ws.send(JSON.stringify(request));
};

export {
  getLastMessages,
  createMessageWebSocket,
  sendMessageWS,
  sendDeleteRequest,
};

