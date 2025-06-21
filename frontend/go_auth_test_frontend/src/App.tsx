import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import Home, { BACKEND_URL } from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import { useEffect } from "react";
import { useDispatch } from "react-redux";
import { setUser, type User } from "./store/store";
import VerifyEmail from "./pages/VerifyEmail";

function App() {
  const dispatch = useDispatch();
  const token = localStorage.getItem("JWT");
  
  const fetchUser = async () => {
    if (token) {
      console.log("Fetching User");
      try {
        const response = await fetch(BACKEND_URL + "/api/profile", {
          headers: {
            "Content-Type": "application/json",
            Authorization: "Bearer " + token,
          },
        });
        if (!response.ok) {
          console.log(response);
        }

        const json = (await response.json()) as User;
        console.log(json);
        dispatch(setUser(json));
      } catch (error) {
        console.log(error);
      }
    }
  };

  useEffect(() => {
    fetchUser();
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/account-created" element={<VerifyEmail />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
