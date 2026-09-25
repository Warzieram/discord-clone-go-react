import { useState, type FormEvent } from "react";
import { NavLink } from "react-router-dom";
import BurgerIcon from "./BurgerIcon";
import RoomCreationForm from "./RoomCreationForm";
import UserAvatar from "./UserAvatar";
import { useSelector } from "react-redux";
import type { RootState } from "../store/store";
import type { Room } from "../types";

type RoomListProps = {
  rooms: Array<Room>;
  onCreateRoom: (name: string) => Promise<void>;
};

const MAX_ROOM_NAME_LENGTH = 20;

const RoomList = ({ rooms, onCreateRoom }: RoomListProps) => {
  const [name, setName] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [hidden, setHidden] = useState<boolean>(false);
  const [createFormHidden, setCreateFormHidden] = useState<boolean>(true);
  const username = useSelector((state: RootState) => state.user.user?.username);

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
    setCreateFormHidden(true);
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
          <li>
            <button
              id="create-room-button"
              onClick={() => setCreateFormHidden(false)}
            >
              +
            </button>
          </li>
        </ul>
        {!createFormHidden ? (
          <RoomCreationForm
            onSubmit={handleSubmit}
            onChange={(e) => setName(e.target.value)}
            onCancel={() => setCreateFormHidden(true)}
            error={error}
            name={name}
            maxLength={MAX_ROOM_NAME_LENGTH}
          />
        ) : (
          ""
        )}
        <NavLink to="/">
          <div id="profile-section-sidebar">
            <UserAvatar username={username || ""} />
            <p>{username}</p>
          </div>
        </NavLink>
      </nav>
    );
  } else {
    return <BurgerIcon onClick={() => setHidden(false)} />;
  }
};

export default RoomList;
