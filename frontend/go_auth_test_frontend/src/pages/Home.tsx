import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { clearToken, logout, type RootState } from "../store/store";
import { parseISO } from "date-fns";
import RedirectionButton from "../components/RedirectionButton";

export const BACKEND_URL = "http://localhost:8080";

export type Room = {
  id: number;
  name: string;
  creator_id: number;
};

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
  }, [user, token]);

  if (!user)
    return (
      <>
        <RedirectionButton to="/login" variation="dark">
          Login
        </RedirectionButton>
        <RedirectionButton to="/register" variation="light">
          Register
        </RedirectionButton>
      </>
    );

  return (
    <>
      <h2>Profil utilisateur</h2>
      <p>Email: {user.email}</p>
      <p>Created on: {creationDate}</p>
      <button onClick={handleLogout}>Logout</button>
    </>
  );
};

export default Home;
