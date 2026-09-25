import type { RefObject } from "react";
import type { Message } from "./MessageCard";

type MessageMenuProps = {
  menuRef: RefObject<HTMLDivElement | null>;
  menuPosition: { x: number; y: number };
  handleCopyMessage: () => void;
  handleDeleteMessage: () => void;
  currentUser: string;
  message: Message;
};

const MessageMenu = ({
  menuRef,
  menuPosition,
  handleCopyMessage,
  handleDeleteMessage,
  currentUser,
  message,
}: MessageMenuProps) => {
  return (
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
  );
};

export default MessageMenu;
