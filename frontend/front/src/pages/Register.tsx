import { useDispatch, useSelector } from "react-redux";
import { setUser, type RootState } from "../store/store";
import { useEffect, useState } from "react";
import RegisterForm from "../components/RegisterForm";
import { useNavigate } from "react-router-dom";
import { regsiter } from "../services/authService";
import type { AuthApiResponse, RegisterFormReturn } from "../types";

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

  const handleRegister = async ({
    email,
    password,
    username,
  }: RegisterFormReturn) => {
    try {
      const response = await regsiter(email, password, username);
      if (!response.ok) {
        throw new Error(await response.text());
      }

      const json = (await response.json()) as AuthApiResponse;
      console.log(json);

      dispatch(setUser(json.user));
      navigate("/account-created");
    } catch (err) {
      console.log(err);
      const error = err as Error;
      setError(error.message);
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
