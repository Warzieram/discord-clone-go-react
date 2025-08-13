import {
  useEffect,
  useRef,
  useState,
  type ChangeEvent,
  type MouseEventHandler,
} from "react";
import { useSelector } from "react-redux";
import { useNavigate, useParams } from "react-router-dom";
import type { RootState } from "../store/store";
import type { Message } from "../components/MessageCard";
import MessageCard from "../components/MessageCard";
import { BACKEND_URL } from "./Home";

type BroadcastedMessage = {
  command_type: string;
  data: Message | number;
};

const ChatRoom = () => {
  const [lastMessage, setLastMessage] = useState<Message>();
  const [messages, setMessages] = useState<Array<Message>>([]);
  const [input, setInput] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const token = useSelector((state: RootState) => state.token.token);
  const navigate = useNavigate();
  const username = useSelector((state: RootState) => state.user.user?.username);
  const ws = useRef<WebSocket | null>(null);
  const params = useParams()
  const id = params.id

  const handleType = (e: ChangeEvent<HTMLInputElement>) => {
    setInput(e.target.value);
  };

  useEffect(() => {
    if (!token) {
      navigate("/login");
    }

    const deleteMessage = (id: number) => {
      setMessages(prev => prev.filter((m) => m.id !== id));
    };

    const retrieveMessages = async () => {
      try {
        const res = await fetch(
          BACKEND_URL + "/api/messages?limit=10&offset=0&room="+id,
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
        console.log(retrievedMessages);
        if (!retrievedMessages) {
          setMessages([])
        }
        else{
          setMessages(retrievedMessages.reverse());
        }
      } catch (error) {
        console.error(error);
        const err = error as Error;
        setError(err.message);
      }
    };

    retrieveMessages()

    ws.current = new WebSocket("ws://localhost:8080/api/message", [
      `auth.${token}`,
    ]);
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
        deleteMessage(data.data as number);
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
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token, navigate]);

  useEffect(() => {
    if (lastMessage) {
      setMessages((old) => [...old, lastMessage]);
    }
  }, [lastMessage]);

  const sendMessage: MouseEventHandler<HTMLButtonElement> = (e) => {
    e.preventDefault();
    if (input && ws.current && ws.current.readyState == WebSocket.OPEN) {
      const request = {
        command_type: "SEND",
        data: {
          content: input,
          room_id: parseInt(id || "-1")
        },
      };
      ws.current.send(JSON.stringify(request));
      setInput("");
    }
  };

  const sendDeleteRequest = (id: number) => {
    const request = {
      command_type: "REMOVE",
      data: id,
    };
    ws.current?.send(JSON.stringify(request));
  };

  return (
    <>
    <h1>{}</h1>
      {messages.map((message: Message, id: number) => (
        <MessageCard
          message={message}
          key={id}
          currentUser={username}
          onDeleteMessage={sendDeleteRequest}
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
