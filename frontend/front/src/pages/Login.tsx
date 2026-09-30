import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { setToken, setUser, type RootState } from "../store/store";
import LoginForm from "../components/LoginForm";
import { Link, useNavigate } from "react-router-dom";
import { login } from "../services/authService";
import type { AuthApiResponse, LoginFormReturn } from "../types";

const Login = () => {
  const [error, setError] = useState<string>();
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const token = useSelector((state: RootState) => state.token.token);

  useEffect(() => {
    if (token) {
      navigate("/");
    }
  }, [token, navigate]);

  const handleLogin = async ({ email, password }: LoginFormReturn) => {
    try {
      const response = await login(email, password);

      if (!response.ok) {
        throw new Error(await response.text());
      }

      const json = (await response.json()) as AuthApiResponse;

      dispatch(setToken(json.token));
      dispatch(setUser(json.user));
      navigate("/");
    } catch (err) {
      const error = err as Error;
      setError(error.message);
    }
  };

  return (
    <div className="auth-page">
      <div className="form-card">
        <h2 className="auth-title">Content de te revoir !</h2>
        <p className="auth-subtitle">
          On est ravis de te revoir parmi nous
        </p>

        <LoginForm callback={handleLogin} />

        {error && <p className="auth-error">{error}</p>}

        <p className="auth-switch">
          Pas encore de compte ? <Link to="/register">S'inscrire</Link>
        </p>
      </div>
    </div>
  );
};

export default Login;
