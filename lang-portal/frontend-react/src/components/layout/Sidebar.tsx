
import { cn } from "@/lib/utils";
import { Menu } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { Link, useLocation } from "react-router-dom";

const navItems = [
  { href: "/dashboard", label: "Dashboard", icon: "home" },
  { href: "/study-activities", label: "Study Activities", icon: "book" },
  { href: "/words", label: "Words", icon: "file-text" },
  { href: "/groups", label: "Word Groups", icon: "folder" },
  { href: "/sessions", label: "Sessions", icon: "list" },
  { href: "/settings", label: "Settings", icon: "settings" },
];

export const Sidebar = () => {
  const [collapsed, setCollapsed] = useState(false);
  const location = useLocation();

  return (
    <div
      className={cn(
        "h-screen fixed left-0 top-0 z-40 flex flex-col border-r bg-background transition-smooth",
        collapsed ? "w-16" : "w-64"
      )}
    >
      <div className="flex h-14 items-center border-b px-4">
        <Button
          variant="ghost"
          size="icon"
          className="mr-2"
          onClick={() => setCollapsed(!collapsed)}
        >
          <Menu className="h-4 w-4" />
        </Button>
        {!collapsed && (
          <span className="font-semibold text-lg">LangPortal</span>
        )}
      </div>

      <nav className="flex-1 space-y-1 p-2">
        {navItems.map((item) => (
          <Link
            key={item.href}
            to={item.href}
            className={cn(
              "flex items-center rounded-lg px-3 py-2 text-sm transition-smooth",
              location.pathname === item.href
                ? "bg-primary text-primary-foreground"
                : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
            )}
          >
            <i className={`lucide-${item.icon} mr-2 h-4 w-4`} />
            {!collapsed && <span>{item.label}</span>}
          </Link>
        ))}
      </nav>
    </div>
  );
};
