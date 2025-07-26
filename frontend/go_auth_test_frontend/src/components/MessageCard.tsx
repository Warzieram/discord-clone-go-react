export type Message = {
  sender: string;
  created_at: string;
  content: string;
};

type MessageCardProps = {
  message: Message;
};

const MessageCard = ({ message }: MessageCardProps) => {
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

  const getAvatarColor = (username: string) => {
    const colors = [
      '#7289da', '#99aab5', '#f04747', '#faa61a', 
      '#43b581', '#9266cc', '#e91e63', '#00bcd4',
      '#4caf50', '#ff9800', '#795548', '#607d8b'
    ];
    let hash = 0;
    for (let i = 0; i < username.length; i++) {
      hash = username.charCodeAt(i) + ((hash << 5) - hash);
    }
    return colors[Math.abs(hash) % colors.length];
  };

  const avatarColor = getAvatarColor(message.sender);
  const userInitial = message.sender.charAt(0).toUpperCase();

  return (
    <div className="discord-message" title={formattedFullDate}>
      <div className="message-avatar" style={{ backgroundColor: avatarColor }}>
        {userInitial}
      </div>
      <div className="message-content-wrapper">
        <div className="message-header-inline">
          <span className="message-username">{message.sender}</span>
          <span className="message-timestamp">{formattedDate}</span>
        </div>
        <div className="message-text">
          {message.content}
        </div>
      </div>
    </div>
  );
};

export default MessageCard;
