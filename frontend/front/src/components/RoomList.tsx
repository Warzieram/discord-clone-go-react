import { useState, type FormEvent } from "react";
import { NavLink } from "react-router-dom";
import type { Room } from "../services/roomService";

const MAX_ROOM_NAME_LENGTH = 20;

type RoomListProps = {
  rooms: Array<Room>;
  onCreateRoom: (name: string) => Promise<void>;
};

const RoomList = ({ rooms, onCreateRoom }: RoomListProps) => {
  const [name, setName] = useState<string>("");
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const trimmed = name.trim();
    if (!trimmed) return;

    try {
      await onCreateRoom(trimmed);
      setName("");
      setError(null);
    } catch (error) {
      setError((error as Error).message);
    }
  };

  return (
    <nav className="room-sidebar">
      <div className="room-sidebar-title">Rooms</div>
      <ul className="room-list">
        {rooms.map((room) => (
          <li key={room.id}>
            <NavLink
              to={`/chatroom/${room.id}`}
              className={({ isActive }) =>
                isActive ? "room-link active" : "room-link"
              }
            >
              # {room.name}
            </NavLink>
          </li>
        ))}
      </ul>
      <form className="room-create-form" onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="New room"
          value={name}
          maxLength={MAX_ROOM_NAME_LENGTH}
          onChange={(e) => setName(e.target.value)}
        />
        <button type="submit" disabled={!name.trim()}>
          Create
        </button>
      </form>
      {error && <div className="room-error">{error}</div>}
    </nav>
  );
};

export default RoomList;
