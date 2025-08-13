import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import { type RootState } from "../store/store";
import CreateRoomForm, {
  type CreateRoomFormReturn,
} from "../components/CreateRoomForm";
import { BACKEND_URL, type Room } from "./Home";

type CreateRoomRequestBody = {
  name: string;
};

const CreateRoom = () => {
  const [error, setError] = useState<string>();
  const [success, setSuccess] = useState<string>();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const token = useSelector((state: RootState) => state.token.token);
  const user = useSelector((state: RootState) => state.user.user)
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) {
      navigate("/login");
    }
  }, [token, navigate]);

  const handleCreateRoom = async ({ name }: CreateRoomFormReturn) => {
    setIsLoading(true);
    setError("");
    setSuccess("");

    if (!user) {
      return
    }
    try {
      const body: CreateRoomRequestBody = {
        name,
      };
      const response = await fetch(BACKEND_URL + "/api/room", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + token,
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        throw new Error(await response.text());
      }

      const newRoom = (await response.json()) as Room;
      setSuccess(`Room "${name}" created successfully!`);

      // Redirect to the new room after a short delay
      setTimeout(() => {
        navigate(`/chatroom/${newRoom.id}`);
      }, 1500);
    } catch (err) {
      console.error(err);
      const error = err as Error;
      setError(error.message || "Failed to create room. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  if (!token) return null;

  return (
    <div className="page-container">
      <div className="page-content">
        <h2>Create New Room</h2>
        <p className="page-description">
          Create a new chat room for your community. Choose a unique name that
          represents the purpose of your room.
        </p>

        <CreateRoomForm callback={handleCreateRoom} isLoading={isLoading} />

        {error && (
          <div className="message error-message">
            <span className="message-icon">⚠️</span>
            {error}
          </div>
        )}

        {success && (
          <div className="message success-message">
            <span className="message-icon">✅</span>
            {success}
          </div>
        )}

        <div className="page-actions">
          <button
            onClick={() => navigate("/")}
            className="btn-secondary"
            disabled={isLoading}
          >
            Back to Home
          </button>
        </div>
      </div>
    </div>
  );
};

export default CreateRoom;
