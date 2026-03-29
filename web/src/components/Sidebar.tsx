"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

const navItems = [
  { href: "/home", icon: "home", label: "Home" },
  { href: "/search", icon: "search", label: "Search" },
  { href: "/notifications", icon: "bell", label: "Notifications" },
  { href: "/messages", icon: "mail", label: "Messages" },
  { href: "/bookmarks", icon: "bookmark", label: "Bookmarks" },
  { href: "/profile", icon: "user", label: "Profile" },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="fixed left-0 top-0 h-full w-[275px] border-r border-[var(--qube-border)] flex flex-col p-4">
      {/* Logo */}
      <Link href="/home" className="text-3xl font-bold px-4 py-3">
        <span className="text-[var(--qube-accent)]">Q</span>ube
      </Link>

      {/* Nav */}
      <nav className="mt-4 flex flex-col gap-1">
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-4 px-4 py-3 rounded-full text-xl transition-colors",
                isActive ? "font-bold" : "font-normal",
                "hover:bg-[var(--qube-surface-hover)]"
              )}
            >
              <NavIcon name={item.icon} active={isActive} />
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>

      {/* Post button */}
      <Link
        href="/compose"
        className="mt-4 w-full bg-[var(--qube-primary)] hover:bg-[var(--qube-primary-dark)] text-white font-bold text-lg py-3 rounded-full text-center transition-colors"
      >
        Post
      </Link>

      {/* User at bottom */}
      <div className="mt-auto p-4">
        {/* Populated after auth */}
      </div>
    </aside>
  );
}

function NavIcon({ name, active }: { name: string; active: boolean }) {
  const icons: Record<string, React.ReactNode> = {
    home: (
      <svg className="w-7 h-7" fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={active ? "0" : "2"} viewBox="0 0 24 24">
        <path d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
      </svg>
    ),
    search: (
      <svg className="w-7 h-7" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
        <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    ),
    bell: (
      <svg className="w-7 h-7" fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={active ? "0" : "2"} viewBox="0 0 24 24">
        <path d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
      </svg>
    ),
    mail: (
      <svg className="w-7 h-7" fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={active ? "0" : "2"} viewBox="0 0 24 24">
        <path d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
      </svg>
    ),
    bookmark: (
      <svg className="w-7 h-7" fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={active ? "0" : "2"} viewBox="0 0 24 24">
        <path d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
      </svg>
    ),
    user: (
      <svg className="w-7 h-7" fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={active ? "0" : "2"} viewBox="0 0 24 24">
        <path d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
      </svg>
    ),
  };
  return icons[name] || null;
}
