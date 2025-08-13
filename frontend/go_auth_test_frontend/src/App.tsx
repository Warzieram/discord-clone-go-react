import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import Home, { BACKEND_URL } from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import { useEffect } from "react";
import { useDispatch } from "react-redux";
import { clearToken, setUser, type User } from "./store/store";
import VerifyEmail from "./pages/VerifyEmail";
import ChatRoom from "./pages/ChatRoom";
import Layout from "./components/Layout";
import CreateRoom from "./pages/CreateRoom";

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
          dispatch(clearToken())
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
  }, [token]);

  useEffect(() => {
    fetchUser();
  }, []);

  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/account-created" element={<VerifyEmail />} />
          <Route path="/chatroom/:id" element={<ChatRoom />} />
          <Route path="/newroom" element={<CreateRoom />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}

export default App;
