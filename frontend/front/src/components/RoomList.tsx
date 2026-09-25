import { useState, type FormEvent } from "react";
import { NavLink } from "react-router-dom";
import type { Room } from "../services/roomService";
import BurgerIcon from "./BurgerIcon";
import RoomCreationForm from "./RoomCreationForm";

const MAX_ROOM_NAME_LENGTH = 20;

type RoomListProps = {
  rooms: Array<Room>;
  onCreateRoom: (name: string) => Promise<void>;
};

const RoomList = ({ rooms, onCreateRoom }: RoomListProps) => {
  const [name, setName] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [hidden, setHidden] = useState<boolean>(false);

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

  if (!hidden) {
    return (
      <nav className="room-sidebar">
        <BurgerIcon onClick={() => setHidden(true)} />
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
        <RoomCreationForm
          onSubmit={handleSubmit}
          onChange={(e) => setName(e.target.value)}
          error={error}
          name={name}
          maxLength={MAX_ROOM_NAME_LENGTH}
        />
      </nav>
    );
  } else {
    return <BurgerIcon onClick={() => setHidden(false)} />;
  }
};

export default RoomList;
