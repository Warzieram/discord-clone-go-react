import { useState, useRef, useEffect } from "react";
import UserAvatar from "./UserAvatar";
import MessageMenu from "./MessageMenu";

export type Message = {
  id?: number;
  sender: string;
  created_at: string;
  content: string;
  room_id?: number;
};

type MessageCardProps = {
  message: Message;
  onDeleteMessage?: (messageId: number) => void;
  currentUser?: string;
};

const MessageCard = ({
  message,
  onDeleteMessage,
  currentUser,
}: MessageCardProps) => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ x: 0, y: 0 });
  const menuRef = useRef<HTMLDivElement>(null);

  const date = new Date(message.created_at);

  const formattedDate = date.toLocaleString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
  });

  const formattedFullDate = date.toLocaleString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

  const handleDotsClick = (event: React.MouseEvent) => {
    event.stopPropagation();
    const rect = event.currentTarget.getBoundingClientRect();
    setMenuPosition({
      x: rect.right - 150,
      y: rect.bottom + 5,
    });
    setIsMenuOpen(!isMenuOpen);
  };
  const handleDeleteMessage = () => {
    if (onDeleteMessage && message.id) {
      onDeleteMessage(message.id);
    }
    setIsMenuOpen(false);
  };

  const handleCopyMessage = () => {
    navigator.clipboard.writeText(message.content);
    setIsMenuOpen(false);
  };

  // Click outside to close menu
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false);
      }
    };

    if (isMenuOpen) {
      document.addEventListener("mousedown", handleClickOutside);
      return () =>
        document.removeEventListener("mousedown", handleClickOutside);
    }
  }, [isMenuOpen]);

  return (
    <div className="discord-message" title={formattedFullDate}>
      <UserAvatar username={message.sender} />
      <div className="message-content-wrapper">
        <div className="message-header-inline">
          <span className="message-username">{message.sender}</span>
          <span className="message-timestamp">{formattedDate}</span>
          <span className="message-three-dots" onClick={handleDotsClick}>
            ⋯
          </span>
        </div>
        <div className="message-text">{message.content}</div>
      </div>

      {isMenuOpen && (
        <MessageMenu
          message={message}
          handleDeleteMessage={handleDeleteMessage}
          menuPosition={menuPosition}
          handleCopyMessage={handleCopyMessage}
          menuRef={menuRef}
          currentUser={currentUser || ""}
        />
      )}
    </div>
  );
};

export default MessageCard;
