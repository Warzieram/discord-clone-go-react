import { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { Link, useLocation } from "react-router-dom";
import { type RootState } from "../store/store";
import { BACKEND_URL, type Room } from "../pages/Home";

type SidebarProps = {
  isOpen: boolean;
  onClose: () => void;
};

const Sidebar = ({ isOpen, onClose }: SidebarProps) => {
  const user = useSelector((state: RootState) => state.user.user);
  const token = useSelector((state: RootState) => state.token.token);
  const [rooms, setRooms] = useState<Array<Room>>([]);
  const location = useLocation();

  useEffect(() => {
    if (user && token) {
      const fetchRooms = async () => {
        try {
          const res = await fetch(BACKEND_URL + "/api/rooms", {
            headers: {
              "Content-Type": "application/json",
              Authorization: "Bearer " + token,
            },
          });
          if (!res.ok) {
            throw new Error(await res.text());
          }
          const retrievedRooms = (await res.json()) as Array<Room>;
          setRooms(retrievedRooms);
        } catch (error) {
          const err = error as Error;
          console.log(err);
        }
      };
      fetchRooms();
    }
  }, [user, token]);

  if (!user) return null;

  return (
    <>
      {/* Overlay for mobile */}
      {isOpen && <div className="sidebar-overlay" onClick={onClose} />}
      
      {/* Sidebar */}
      <div className={`sidebar ${isOpen ? 'sidebar--open' : ''}`}>
        <div className="sidebar__header">
          <h3>Chat Rooms</h3>
        </div>
        
        <div className="sidebar__content">
          <nav className="sidebar__nav">
            <Link
              to="/"
              className={`sidebar__link ${location.pathname === '/' ? 'sidebar__link--active' : ''}`}
              onClick={onClose}
            >
              <span className="sidebar__link-icon">🏠</span>
              Home
            </Link>
            {rooms.map((room) => {
              const isActive = location.pathname === `/chatroom/${room.id}`;
              return (
                <Link
                  key={room.id}
                  to={`/chatroom/${room.id}`}
                  className={`sidebar__link ${isActive ? 'sidebar__link--active' : ''}`}
                  onClick={onClose}
                >
                  <span className="sidebar__link-icon">#</span>
                  {room.name}
                </Link>
              );
            })}
          </nav>
        </div>
      </div>
    </>
  );
};

export default Sidebar;
