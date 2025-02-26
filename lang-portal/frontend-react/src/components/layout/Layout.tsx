
import { Sidebar } from "./Sidebar";
import { Outlet } from "react-router-dom";

export const Layout = () => {
  return (
    <div className="min-h-screen bg-background">
      <Sidebar />
      <main className="pl-64 transition-smooth">
        <div className="container py-8">
          <Outlet />
        </div>
      </main>
    </div>
  );
};
