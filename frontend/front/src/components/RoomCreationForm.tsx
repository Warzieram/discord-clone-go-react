import type { ChangeEventHandler, FormEvent } from "react";

type RoomCreationFormProps = {
  onSubmit: (e: FormEvent<HTMLFormElement>) => Promise<void>;
  onCancel: () => void;
  name: string;
  maxLength: number;
  onChange: ChangeEventHandler<HTMLInputElement>;
  error: string | null;
};

const RoomCreationForm = ({
  onSubmit,
  onChange,
  onCancel,
  name,
  maxLength,
  error,
}: RoomCreationFormProps) => {
  return (
    <>
      <form className="room-create-form" onSubmit={onSubmit}>
        <input
          type="text"
          placeholder="New room"
          value={name}
          maxLength={maxLength}
          onChange={onChange}
        />
        <div id="room-create-form-button-section">
          <button type="submit" disabled={!name.trim()}>
            Create
          </button>
          <button className="btn-secondary" onClick={onCancel}>
            Cancel
          </button>
        </div>
      </form>
      {error && <div className="room-error">{error}</div>}
    </>
  );
};

export default RoomCreationForm;
