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
  const [passwordVisible, setPasswordVisible] = useState<boolean>(false)

  const handleChange = (
    event: ChangeEvent<HTMLInputElement>,
    onChangeCallback: Dispatch<SetStateAction<string>>,
  ) => {
    event.preventDefault();
    onChangeCallback(event.target.value);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    callback({ email, password });
  };

  return (
    <div className="form-card">
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            name="email"
            id="email"
            value={email}
            onChange={(e) => handleChange(e, setEmail)}
            placeholder="your@email.com"
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            type={passwordVisible ? "text" : "password"}
            name="password"
            id="password"
            value={password}
            onChange={(e) => handleChange(e, setPassword)}
            placeholder="Enter your password"
            required
          />
          <div className="password-toggle">
            <input
              type="checkbox"
              name="passwordVisible"
              id="passwordVisible"
              checked={passwordVisible}
              onChange={() => setPasswordVisible(!passwordVisible)}
            />
            <label htmlFor="passwordVisible" className="checkbox-label">
              Show Password
            </label>
          </div>
        </div>
        <button type="submit" className="btn-primary">
          Sign In
        </button>
      </form>
    </div>
  );
}

export default LoginForm;
