import { useEffect, useState } from "react";

type UserAvatarProps = {
  username: string;
  onClick?: () => void;
};

const UserAvatar = ({ username, onClick }: UserAvatarProps) => {
  const [avatarColor, setAvatarColor] = useState<string>();
  const [initial, setInitial] = useState<string>();
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

    const color = getAvatarColor(username);
    setAvatarColor(color);
    const userInitial = username.charAt(0).toUpperCase();
    setInitial(userInitial);
  }, [username]);

  return (
    <div
      onClick={onClick}
      className="message-avatar"
      style={{ backgroundColor: avatarColor }}
    >
      {initial}
    </div>
  );
};

export default UserAvatar;
