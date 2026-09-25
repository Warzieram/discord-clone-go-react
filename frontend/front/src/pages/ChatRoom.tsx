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
import RoomList from "../components/RoomList";
import {
  createMessageWebSocket,
  getLastMessages,
  sendDeleteRequest,
  sendMessageWS,
  type BroadcastedMessage,
} from "../services/messageService";
import { createRoom, getRooms, type Room } from "../services/roomService";

const ChatRoom = () => {
  const [messages, setMessages] = useState<Array<Message>>([]);
  const [rooms, setRooms] = useState<Array<Room>>([]);
  const [roomsLoaded, setRoomsLoaded] = useState<boolean>(false);
  const [input, setInput] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const token = useSelector((state: RootState) => state.token.token);
  const navigate = useNavigate();
  const username = useSelector((state: RootState) => state.user.user?.username);
  const ws = useRef<WebSocket | null>(null);

  const { roomId: roomIdParam } = useParams();
  const roomId = roomIdParam ? Number(roomIdParam) : undefined;
  const currentRoom = rooms.find((room) => room.id === roomId);

  const handleType = (e: ChangeEvent<HTMLInputElement>) => {
    setInput(e.target.value);
  };

  // Load the room list once
  useEffect(() => {
    if (!token) {
      navigate("/login");
      return;
    }

    const retrieveRooms = async () => {
      try {
        const res = await getRooms(token);
        if (!res.ok) {
          throw new Error(await res.text());
        }
        setRooms((await res.json()) as Array<Room>);
      } catch (error) {
        console.error(error);
        setError((error as Error).message);
      } finally {
        setRoomsLoaded(true);
      }
    };

    retrieveRooms();
  }, [token, navigate]);

  // Without a room in the URL, open the first one
  useEffect(() => {
    if (roomIdParam === undefined && rooms.length > 0) {
      navigate(`/chatroom/${rooms[0].id}`, { replace: true });
    }
  }, [roomIdParam, rooms, navigate]);

  // (Re)connect to the selected room
  useEffect(() => {
    if (!token || roomId === undefined || !Number.isInteger(roomId)) {
      return;
    }

    // Set when this effect is cleaned up (room change / unmount), so late
    // responses and the resulting close event don't touch the new room's state
    let stale = false;
    setMessages([]);
    setError(null);

    const deleteMessageFromList = (id: number) => {
      setMessages((prev) => prev.filter((m) => m.id !== id));
    };

    const retrieveMessages = async () => {
      try {
        const res = await getLastMessages(token, roomId, 10, 0);
        if (!res.ok) {
          throw new Error(await res.text());
        }
        const retrievedMessages = (await res.json()) as Array<Message>;
        if (!stale && retrievedMessages) {
          // The API returns newest first; display oldest first. Keep any
          // message that arrived over the socket while this was loading.
          setMessages((live) => [
            ...[...retrievedMessages].reverse(),
            ...live.filter(
              (m) => !retrievedMessages.some((r) => r.id === m.id),
            ),
          ]);
        }
      } catch (error) {
        console.error(error);
        if (!stale) {
          setError((error as Error).message);
        }
      }
    };

    retrieveMessages();

    const socket = createMessageWebSocket(token, roomId);
    ws.current = socket;

    socket.addEventListener("open", () => {
      console.log("WS connection established to room", roomId);
      setError("");
    });

    socket.addEventListener("message", (event) => {
      const data: BroadcastedMessage = JSON.parse(event.data);

      if (data.command_type === "REMOVE") {
        deleteMessageFromList(data.data as number);
      } else if (data.command_type === "SEND") {
        setMessages((old) => [...old, data.data as Message]);
      }
    });

    socket.addEventListener("close", () => {
      if (stale) return;
      setError("You got disconnected, please refresh the page");
      console.log("Closed ws connexion");
    });

    return (): void => {
      stale = true;
      socket.close();
    };
  }, [token, roomId]);

  const sendMessage: MouseEventHandler<HTMLButtonElement> = (e) => {
    e.preventDefault();
    if (input && ws.current && ws.current.readyState == WebSocket.OPEN) {
      sendMessageWS(ws.current, input);
      setInput("");
    }
  };

  const handleCreateRoom = async (name: string) => {
    const res = await createRoom(token || "", name);
    if (!res.ok) {
      throw new Error(await res.text());
    }
    const room = (await res.json()) as Room;
    setRooms((prev) => [...prev, room]);
    navigate(`/chatroom/${room.id}`);
  };

  const renderChat = () => {
    if (!roomsLoaded) {
      return null;
    }
    if (rooms.length === 0) {
      return <div className="chat-placeholder">Create a room to start chatting</div>;
    }
    if (!currentRoom) {
      return <div className="chat-placeholder">Room not found</div>;
    }

    return (
      <>
        <div className="chat-header"># {currentRoom.name}</div>
        <div className="chat-messages">
          {messages.map((message: Message, id: number) => (
            <MessageCard
              message={message}
              key={id}
              currentUser={username}
              onDeleteMessage={(id) => sendDeleteRequest(id, ws.current)}
            />
          ))}
        </div>
        <form className="message-form">
          <div className="message-input-form">
            <input
              className="message-input"
              type="text"
              onChange={handleType}
              value={input}
              placeholder={`Message #${currentRoom.name}`}
              autoFocus={true}
            ></input>
            <button className="send-button" type="submit" onClick={sendMessage}>
              Send
            </button>
          </div>
        </form>
      </>
    );
  };

  return (
    <div className="chat-layout">
      <RoomList rooms={rooms} onCreateRoom={handleCreateRoom} />
      <main className="chat-main">
        {renderChat()}
        {error && <div className="chat-error"> {error} </div>}
      </main>
    </div>
  );
};

export default ChatRoom;
