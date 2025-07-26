import {
  useEffect,
  useRef,
  useState,
  type ChangeEvent,
  type MouseEventHandler,
} from "react";
import { useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import type { RootState } from "../store/store";
import type { Message } from "../components/MessageCard";
import MessageCard from "../components/MessageCard";
import { BACKEND_URL } from "./Home";

const ChatRoom = () => {
  const [lastMessage, setLastMessage] = useState<Message>();
  const [messages, setMessages] = useState<Array<Message>>([]);
  const [input, setInput] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const token = useSelector((state: RootState) => state.token.token);
  const navigate = useNavigate();
  const ws = useRef<WebSocket | null>(null);

  const handleType = (e: ChangeEvent<HTMLInputElement>) => {
    setInput(e.target.value);
  };

  useEffect(() => {
    if (!token) {
      navigate("/login");
    }

    const retrieveMessages = async () => {
      try {
        const res = await fetch(
          BACKEND_URL + "/api/messages?limit=10&offset=0",
          {
            headers: {
              "Content-Type": "application/json",
              Authorization: "Bearer " + token,
            },
          },
        );
        if (!res.ok) {
          throw new Error(await res.text());
        }
        const retrievedMessages = (await res.json()) as Array<Message>;
        setMessages(retrievedMessages)
        console.log(retrievedMessages);
      } catch (error) {
        console.error(error);
        const err = error as Error;
        setError(err.message);
      }
    };

    retrieveMessages();

    ws.current = new WebSocket("ws://localhost:8080/api/message", [
      `auth.${token}`,
    ]);
    console.log(ws.current);
    ws.current.addEventListener("open", () => {
      console.log("WS connection established");
      setError("");
    });

    ws.current.addEventListener("message", (event) => {
      console.log(event.data);
      setLastMessage(JSON.parse(event.data));
    });

    ws.current.addEventListener("close", () => {
      setError("You got disconnected, please refresh the page");
      console.log("Closed ws connexion");
    });

    return (): void => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [token, navigate]);

  useEffect(() => {
    if (lastMessage) {
      setMessages((old) => [...old, lastMessage]);
    }
  }, [lastMessage]);

  const sendMessage: MouseEventHandler<HTMLButtonElement> = (e) => {
    e.preventDefault();
    if (input && ws.current && ws.current.readyState == WebSocket.OPEN) {
      ws.current.send("SEND " + input);
      setInput("");
    }
  };
  return (
    <>
      {messages.map((message: Message, id: number) => (
        <MessageCard message={message} key={id} />
      ))}
      <form className="message-form">
        <div className="message-input-form">
          <input
            className="message-input"
            type="text"
            onChange={handleType}
            value={input}
            autoFocus={true}
          ></input>
          <button type="submit" onClick={sendMessage}>
            Send
          </button>
        </div>
      </form>
      {error && <div> {error} </div>}
    </>
  );
};

export default ChatRoom;
