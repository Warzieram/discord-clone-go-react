import { useState, useRef, useEffect } from "react";

export type Message = {
  id?: number;
  sender: string;
  created_at: string;
  content: string;
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
  const [avatarColor, setAvatarColor] = useState<string>()
  const [initial, setInitial] = useState<string>()
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

  useEffect(() => {
    const getAvatarColor = (username: string) => {
      const colors = [
        "#7289da",
        "#99aab5",
        "#f04747",
        "#faa61a",
        "#43b581",
        "#9266cc",
        "#e91e63",
        "#00bcd4",
        "#4caf50",
        "#ff9800",
        "#795548",
        "#607d8b",
      ];
      let hash = 0;
      for (let i = 0; i < username.length; i++) {
        hash = username.charCodeAt(i) + ((hash << 5) - hash);
      }
      return colors[Math.abs(hash) % colors.length];
    };

    const color = getAvatarColor(message.sender);
    setAvatarColor(color)
    const userInitial = message.sender.charAt(0).toUpperCase();
    setInitial(userInitial)

  }, [message.sender]);

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
      <div className="message-avatar" style={{ backgroundColor: avatarColor }}>
        {initial}
      </div>
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
        <div
          ref={menuRef}
          className="message-context-menu"
          style={{
            position: "fixed",
            left: menuPosition.x,
            top: menuPosition.y,
            zIndex: 1000,
          }}
        >
          <div className="menu-item" onClick={handleCopyMessage}>
            Copy Message
          </div>
          {currentUser === message.sender && (
            <div className="menu-item delete" onClick={handleDeleteMessage}>
              Delete Message
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default MessageCard;
