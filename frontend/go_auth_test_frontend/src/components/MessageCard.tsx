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

  const formattedDate = date.toLocaleString("fr-FR", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

  return (
    <div className="message-container">
      <div className="message-header">
        <p className="message-sender">{message.sender}</p>
      </div>
      <div className="message-body">
        <p className="message-content">{message.content}</p>
      </div>
      <div className="message-footer">
        <p className="message-sent-at">{formattedDate}</p>
      </div>
    </div>
  );
};

export default MessageCard;
