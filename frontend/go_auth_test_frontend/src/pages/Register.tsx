import { useDispatch, useSelector } from "react-redux";
import {  setUser, type RootState, type User } from "../store/store";
import { BACKEND_URL } from "./Home";
import { useEffect, useState } from "react";
import RegisterForm from "../components/RegisterForm";
import { useNavigate } from "react-router-dom";

type RegisterResponseApiType = {
  token: string;
  user: User;
};

type RegisterFormReturn = {
  email: string;
  password: string;
};

const Register = () => {
  const token = useSelector((state: RootState) => state.token.token);
  const dispatch = useDispatch();
  const [error, setError] = useState<string>();
  const navigate = useNavigate();

  useEffect(() => {
    if (token) {
      navigate("/");
    }
  }, [token, navigate]);

  const handleRegister = async ({ email, password }: RegisterFormReturn) => {
    try {
      const response = await fetch(BACKEND_URL + "/api/register", {
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
        throw new Error(await response.text());
      }

      const json = (await response.json()) as RegisterResponseApiType;
      console.log(json);

      dispatch(setUser(json.user));
      navigate("/account-created");
    } catch (err) {
      console.log(err);
      setError(err.message);
    }
  };

  return (
    <>
      <h2>Register</h2>
      <RegisterForm callback={handleRegister} />
      {error}
    </>
  );
};

export default Register;
