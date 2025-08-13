import {
  useState,
  type ChangeEvent,
  type Dispatch,
  type SetStateAction,
} from "react";

type CreateRoomFormProps = {
  callback: (args: CreateRoomFormReturn) => Promise<void>;
  isLoading?: boolean;
};

export type CreateRoomFormReturn = {
  name: string;
};

const CreateRoomForm = ({ callback, isLoading = false }: CreateRoomFormProps) => {
  const [name, setName] = useState<string>("");
  const [nameError, setNameError] = useState<string>("");

  const handleChange = (
    event: ChangeEvent<HTMLInputElement>,
    onChangeCallback: Dispatch<SetStateAction<string>>,
  ) => {
    const value = event.target.value;
    onChangeCallback(value);
    
    // Clear error when user starts typing
    if (nameError) {
      setNameError("");
    }
    
    // Validate max length (15 characters as per schema)
    if (value.length > 15) {
      setNameError("Room name must be 15 characters or less");
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Validation
    if (!name.trim()) {
      setNameError("Room name is required");
      return;
    }
    
    if (name.length > 15) {
      setNameError("Room name must be 15 characters or less");
      return;
    }
    
    callback({ name: name.trim() });
  };

  return (
    <div className="form-card">
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="name">Room Name</label>
          <input
            type="text"
            name="name"
            id="name"
            value={name}
            onChange={(e) => handleChange(e, setName)}
            placeholder="Enter room name..."
            maxLength={15}
            className={nameError ? "input-error" : ""}
            disabled={isLoading}
          />
          {nameError && <span className="error-text">{nameError}</span>}
          <div className="char-counter">
            {name.length}/15 characters
          </div>
        </div>

        <button
          type="submit"
          disabled={isLoading || !name.trim() || name.length > 15}
          className="btn-primary"
        >
          {isLoading ? "Creating..." : "Create Room"}
        </button>
      </form>
    </div>
  );
};

export default CreateRoomForm;