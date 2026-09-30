import {
  useState,
  type ChangeEvent,
  type Dispatch,
  type SetStateAction,
} from "react";
import type { LoginFormReturn } from "../types";

type LoginFormProps = {
  callback: (args: LoginFormReturn) => Promise<void>;
};

function LoginForm({ callback }: LoginFormProps) {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");

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
        callback({ email: email, password: password });
      }}
    >
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
          autoComplete="current-password"
          onChange={(e) => handleChange(e, setPassword)}
          placeholder="••••••••"
        />
      </div>

      <button type="submit" className="auth-submit">
        Se Connecter
      </button>
    </form>
  );
}

export default LoginForm;
