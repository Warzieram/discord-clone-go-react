import {  useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { setToken, setUser, type RootState, type User } from "../store/store";
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
  const token = useSelector((state: RootState) => state.token.token)

  useEffect(() => {
    if (token) {
      navigate("/");
    }
  }, [token, navigate]);

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
      console.log(json.token);
      
      dispatch(setToken(json.token));
      dispatch(setUser(json.user));
      navigate("/");
    } catch (err) {
      console.log(err);
      const error = err as Error
      setError(error.message);
    }
  };

  return (
    <div className="page-container">
      <div className="page-content">
        <h2>Welcome Back</h2>
        <p className="page-description">
          Sign in to your account to continue chatting with your community.
        </p>
        
        <LoginForm callback={handleLogin} />
        
        {error && (
          <div className="message error-message">
            <span className="message-icon">⚠️</span>
            {error}
          </div>
        )}
        
        <div className="page-actions">
          <p style={{ color: '#b9bbbe', fontSize: '14px', marginTop: '16px' }}>
            Don't have an account?{' '}
            <a 
              href="/register" 
              style={{ color: '#5865f2', textDecoration: 'none' }}
              onMouseOver={(e) => e.target.style.textDecoration = 'underline'}
              onMouseOut={(e) => e.target.style.textDecoration = 'none'}
            >
              Sign up here
            </a>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;
