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
    callback({ email, password, username });
  };

  return (
    <div className="form-card">
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">Username</label>
          <input
            type="text"
            name="username"
            id="username"
            value={username}
            onChange={(e) => handleChange(e, setUsername)}
            placeholder="Choose a username"
            required
          />
        </div>
        
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
            placeholder="Create a strong password"
            required
          />
          <div className="password-toggle">
            <input
              type="checkbox"
              name="passwordVisible"
              id="passwordVisibleRegister"
              checked={passwordVisible}
              onChange={() => setPasswordVisible(!passwordVisible)}
            />
            <label htmlFor="passwordVisibleRegister" className="checkbox-label">
              Show Password
            </label>
          </div>
        </div>

        <button type="submit" className="btn-primary">
          Create Account
        </button>
      </form>
    </div>
  );
};

export default RegisterForm;
