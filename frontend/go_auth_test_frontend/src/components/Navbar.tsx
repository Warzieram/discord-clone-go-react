import { useSelector } from "react-redux";
import { type RootState } from "../store/store";

type NavbarProps = {
  onToggleSidebar: () => void;
};

const Navbar = ({ onToggleSidebar }: NavbarProps) => {
  const user = useSelector((state: RootState) => state.user.user);

  if (!user) return null;

  return (
    <nav className="navbar">
      <button
        className="navbar__toggle"
        onClick={onToggleSidebar}
        aria-label="Toggle sidebar"
      >
        <span></span>
        <span></span>
        <span></span>
      </button>
    </nav>
  );
};

export default Navbar;