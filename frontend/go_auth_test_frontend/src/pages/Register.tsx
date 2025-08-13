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
  username: string;
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

  const handleRegister = async ({ email, password, username }: RegisterFormReturn) => {
    try {
      const response = await fetch(BACKEND_URL + "/api/register", {
        method: "post",
        headers: {
          "Content-Type": "aplication/json",
        },
        body: JSON.stringify({
          email,
          password,
          username,
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
      const error = err as Error
      setError(error.message);
    }
  };

  return (
    <div className="page-container">
      <div className="page-content">
        <h2>Join Our Community</h2>
        <p className="page-description">
          Create your account to start chatting and connect with others.
        </p>
        
        <RegisterForm callback={handleRegister} />
        
        {error && (
          <div className="message error-message">
            <span className="message-icon">⚠️</span>
            {error}
          </div>
        )}
        
        <div className="page-actions">
          <p style={{ color: '#b9bbbe', fontSize: '14px', marginTop: '16px' }}>
            Already have an account?{' '}
            <a 
              href="/login" 
              style={{ color: '#5865f2', textDecoration: 'none' }}
              onMouseOver={(e) => e.target.style.textDecoration = 'underline'}
              onMouseOut={(e) => e.target.style.textDecoration = 'none'}
            >
              Sign in here
            </a>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Register;
