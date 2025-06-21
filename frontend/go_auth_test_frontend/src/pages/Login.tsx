import {  useState } from "react";
import { useDispatch } from "react-redux";
import { setToken, setUser, type User } from "../store/store";
import LoginForm, { type LoginFormReturn } from "../components/LoginForm";
import { useNavigate } from "react-router-dom";
import { BACKEND_URL } from "./Home";

type LoginApiResponseType = {
  token: string;
  user: User;
};

const Login = () => {
  const [error, setError] = useState<string>();
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleLogin = async ({ email, password }: LoginFormReturn) => {
    try {
      const response = await fetch(BACKEND_URL + "/api/login", {
        method: "post",
        headers: {
          "Content-Type": "aplication/json",
        },
        body: JSON.stringify({
          email,
          password,
        }),
      });
      if (!response.ok) {
        throw new Error(await response.text())
      }

      const json = (await response.json()) as LoginApiResponseType;
      console.log(json);

      dispatch(setToken(json.token));
      dispatch(setUser(json.user));
      navigate("/");
    } catch (err) {
      console.log(err);
      setError(err.message);
    }
  };

  return (
    <>
      <h2>Log In</h2>
      <LoginForm callback={handleLogin} />
      <p>{ error }</p>
    </>
  );
};

export default Login;
