import { useState } from "react";
import { useSelector } from "react-redux";
import { type RootState } from "../store/store";
import Navbar from "./Navbar";
import Sidebar from "./Sidebar";

type LayoutProps = {
  children: React.ReactNode;
};

const Layout = ({ children }: LayoutProps) => {
  const user = useSelector((state: RootState) => state.user.user);
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const closeSidebar = () => {
    setSidebarOpen(false);
  };

  if (!user) {
    return <div className="app-content">{children}</div>;
  }

  return (
    <div className="app-layout">
      <Navbar onToggleSidebar={toggleSidebar} />
      <Sidebar isOpen={sidebarOpen} onClose={closeSidebar} />
      <div className="app-content">
        {children}
      </div>
    </div>
  );
};

export default Layout;