import {
  useState,
  type ChangeEvent,
  type Dispatch,
  type SetStateAction,
} from "react";
import type { RegisterFormReturn } from "../types";

type RegisterFormProps = {
  callback: (args: RegisterFormReturn) => Promise<void>;
};

const RegisterForm = ({ callback }: RegisterFormProps) => {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [username, setUsername] = useState<string>("");

  const handleChange = (
    event: ChangeEvent<HTMLInputElement>,
    onChangeCallback: Dispatch<SetStateAction<string>>,
  ) => {
    event.preventDefault();
    onChangeCallback(event.target.value);
  };

  return (
    <form
      className="auth-form"
      onSubmit={(e) => {
        e.preventDefault();
        callback({ email: email, password: password, username: username });
      }}
    >
      <div className="form-field">
        <label htmlFor="username">Username</label>
        <input
          type="text"
          name="username"
          id="username"
          autoComplete="username"
          onChange={(e) => handleChange(e, setUsername)}
          placeholder="Example123"
        />
      </div>

      <div className="form-field">
        <label htmlFor="email">Email</label>
        <input
          type="email"
          name="email"
          id="email"
          autoComplete="email"
          onChange={(e) => handleChange(e, setEmail)}
          placeholder="example@thing.com"
        />
      </div>

      <div className="form-field">
        <label htmlFor="password">Mot de passe</label>
        <input
          type="password"
          name="password"
          id="password"
          autoComplete="new-password"
          onChange={(e) => handleChange(e, setPassword)}
          placeholder="••••••••"
        />
      </div>

      <button type="submit" className="auth-submit">
        S'inscrire
      </button>
    </form>
  );
};

export default RegisterForm;
