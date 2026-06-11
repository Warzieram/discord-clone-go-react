import type { Message } from "../components/MessageCard";
import { BACKEND_URL, WS_BACKEND_URL } from "../constants";

export type BroadcastedMessage = {
  command_type: string;
  data: Message | number;
};

const getLastMessages = async (
  token: string,
  limit: number,
  offset: number,
) => {
  const res = await fetch(
    `${BACKEND_URL}/api/messages?limit=${limit}&offset=${offset}`,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + token,
      },
    },
  );

  return res;
};

const createMessageWebSocket = (token: string) => {
  return new WebSocket(`${WS_BACKEND_URL}/api/message`, [`auth.${token}`]);
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

