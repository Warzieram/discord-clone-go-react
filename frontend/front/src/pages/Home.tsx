import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { clearToken, logout, type RootState } from "../store/store";
import { parseISO } from "date-fns";
import RedirectionButton from "../components/RedirectionButton";
import UserAvatar from "../components/UserAvatar";
import { Navigate } from "react-router-dom";

const Home = () => {
  const user = useSelector((state: RootState) => state.user.user);
  const token = useSelector((state: RootState) => state.token.token);
  const dispatch = useDispatch();
  const [creationDate, setCreationDate] = useState<string>();

  const handleLogout = async () => {
    dispatch(clearToken());
    dispatch(logout());
    console.log("From logout: ", user);
  };

  useEffect(() => {
    if (user) {
      if (user.created_at) {
        const formatedDate = parseISO(
          user.created_at || "",
        ).toLocaleDateString();
        setCreationDate(formatedDate);
      }
    }
  }, [user]);

  // `token` is read straight from localStorage when the store is created, so it
  // is available on the very first render and is the reliable "logged in"
  // signal. `user` only arrives once App.tsx's /api/profile fetch resolves, so
  // gating on it here bounced logged-in users to /login on every page load.
  if (!token) {
    return <Navigate to="/login" replace />;
  }

  // Logged in, but the profile request hasn't come back yet.
  if (!user) {
    return <div className="center-wrapper">Loading…</div>;
  }

  return (
    <div className="center-wrapper">
      <UserAvatar username={user.username || ""} />
      <h2>Profil utilisateur</h2>
      <p>Email: {user.email}</p>
      <p>Username: {user.username}</p>
      <p>Created on: {creationDate}</p>
      <div id="room-create-form-button-section">
        <button className="btn-danger" onClick={handleLogout}>
          Logout
        </button>
        <RedirectionButton to="/chatroom" variation="light">
          Chat
        </RedirectionButton>
      </div>
    </div>
  );
};

export default Home;
