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
import {
  createMessageWebSocket,
  getLastMessages,
  sendDeleteRequest,
  sendMessageWS,
  type BroadcastedMessage,
} from "../services/messageService";


const ChatRoom = () => {
  const [lastMessage, setLastMessage] = useState<Message>();
  const [messages, setMessages] = useState<Array<Message>>([]);
  const [input, setInput] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const token = useSelector((state: RootState) => state.token.token);
  const navigate = useNavigate();
  const username = useSelector((state: RootState) => state.user.user?.username);
  const ws = useRef<WebSocket | null>(null);

  const handleType = (e: ChangeEvent<HTMLInputElement>) => {
    setInput(e.target.value);
  };

  useEffect(() => {
    if (!token) {
      navigate("/login");
    }

    const deleteMessageFromList = (id: number) => {
      setMessages((prev) => prev.filter((m) => m.id !== id));
    };

    const retrieveMessages = async () => {
      try {
        const res = await getLastMessages(token || "", 10, 0);
        if (!res.ok) {
          throw new Error(await res.text());
        }
        const retrievedMessages = (await res.json()) as Array<Message>;
        console.log(retrievedMessages);
        if (retrievedMessages) {
          setMessages(retrievedMessages);
        }
      } catch (error) {
        console.error(error);
        const err = error as Error;
        setError(err.message);
      }
    };

    retrieveMessages();

    ws.current = createMessageWebSocket(token || "");

    console.log(ws.current);
    ws.current.addEventListener("open", () => {
      console.log("WS connection established");
      setError("");
    });

    ws.current.addEventListener("message", (event) => {
      const data: BroadcastedMessage = JSON.parse(event.data);
      console.log(data);

      if (data.command_type === "REMOVE") {
        console.log("REMOVING", data.data);
        deleteMessageFromList(data.data as number);
      } else if (data.command_type === "SEND") {
        setLastMessage(data.data as Message);
      }
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
      sendMessageWS(ws.current, input);
      setInput("");
    }
  };

  return (
    <>
      {messages.map((message: Message, id: number) => (
        <MessageCard
          message={message}
          key={id}
          currentUser={username}
          onDeleteMessage={(id) => sendDeleteRequest(id, ws.current)}
        />
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
