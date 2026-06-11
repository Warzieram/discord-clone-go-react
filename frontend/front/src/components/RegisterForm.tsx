import {
  useState,
  type ChangeEvent,
  type Dispatch,
  type SetStateAction,
} from "react";

type RegisterFormProps = {
  callback: (args: RegisterFormReturn) => Promise<void>;
};

export type RegisterFormReturn = {
  email: string;
  password: string;
  username: string;
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
    <div className="form-card">
      <form>
        <label htmlFor="username">Username</label>
        <input
          type="text"
          name="username"
          id="username"
          onChange={(e) => handleChange(e, setUsername)}
          placeholder="Example123"
        />
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
      <button
        onClick={(e) => {
          e.preventDefault()
          callback({ email: email, password: password, username: username });
        }}
          type="submit"
      >
        {" "}
        S'inscrire{" "}
      </button>
      </form>
    </div>
  );
};

export default RegisterForm;
