import {
  useState,
  type ChangeEvent,
  type Dispatch,
  type SetStateAction,
} from "react";

type LoginFormProps = {
  callback: (args: LoginFormReturn) => Promise<void>;
};

export type LoginFormReturn = {
  email: string;
  password: string;
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
    <div className="form-card">
      <form>
        <label htmlFor="email">Email</label>
        <input
          type="text"
          name="email"
          id="email"
          onChange={(e) => handleChange(e, setEmail)}
          placeholder="example@thing.com"
        />
        <label htmlFor="password">Password</label>
        <input
          type="password"
          name="password"
          id="password"
          onChange={(e) => {
            handleChange(e, setPassword);
          }}
          placeholder="example@thing.com"
        />
      </form>
      <button
        onClick={() => {
          callback({ email: email, password: password });
        }}
      >
        {" "}
        Se Connecter{" "}
      </button>
    </div>
  );
}

export default LoginForm;
