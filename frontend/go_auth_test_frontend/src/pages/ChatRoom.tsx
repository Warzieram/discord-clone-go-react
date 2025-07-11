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

const ChatRoom = () => {
  const [lastMessage, setLastMessage] = useState<Message>();
  const [messages, setMessages] = useState<Array<Message>>([]);
  const [input, setInput] = useState<string>("");
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

    ws.current = new WebSocket("ws://192.168.1.151:8080/api/message", [
      `auth.${token}`,
    ]);
    console.log(ws.current);
    ws.current.addEventListener("open", () => {
      console.log("WS connection established");
    });

    ws.current.addEventListener("message", (event) => {
      console.log(event.data);
      setLastMessage(JSON.parse(event.data));
    });

    ws.current.addEventListener("close", () => {
      console.log("Closed ws connexion");
    });

    return (): void => {
      if (ws.current) {
        ws.current.close();
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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
      <form>
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
    </>
  );
};

export default ChatRoom;
