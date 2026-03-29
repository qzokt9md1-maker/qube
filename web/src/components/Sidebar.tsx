"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  { href: "/home", label: "Home", icon: "home" },
  { href: "/search", label: "Search", icon: "search" },
  { href: "/notifications", label: "Notifications", icon: "bell" },
  { href: "/messages", label: "Messages", icon: "mail" },
  { href: "/bookmarks", label: "Bookmarks", icon: "bookmark" },
  { href: "/profile", label: "Profile", icon: "user" },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="fixed left-0 top-0 h-full w-[275px] flex flex-col items-center xl:items-stretch xl:px-3 lg:px-0 py-2 border-r border-[var(--qube-border)] xl:w-[275px] lg:w-[88px] max-lg:hidden">
      {/* Logo */}
      <Link href="/home" className="flex items-center justify-center xl:justify-start h-[50px] w-[50px] xl:w-auto xl:px-4 mb-1 rounded-full hover:bg-[var(--qube-surface-hover)] transition-colors">
        <span className="text-2xl font-black tracking-tight">
          <span className="text-[var(--qube-primary)]">Q</span><span className="xl:inline lg:hidden">ube</span>
        </span>
      </Link>

      {/* Nav */}
      <nav className="flex flex-col gap-0.5 mt-1">
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center justify-center xl:justify-start gap-5 w-[50px] h-[50px] xl:w-auto xl:h-auto xl:px-4 xl:py-3 rounded-full text-[15px] transition-colors hover:bg-[var(--qube-surface-hover)] ${isActive ? "font-bold text-[var(--qube-primary)]" : ""}`}
            >
              <NavIcon name={item.icon} active={isActive} />
              <span className="xl:inline lg:hidden">{item.label}</span>
            </Link>
          );
        })}
      </nav>

      {/* Post button */}
      <button className="mt-4 w-[50px] h-[50px] xl:w-full xl:h-auto xl:py-3 flex items-center justify-center bg-[var(--qube-primary)] hover:bg-[var(--qube-primary-dark)] text-white font-bold rounded-full transition-colors">
        <span className="xl:inline lg:hidden">Post</span>
        <svg className="w-6 h-6 xl:hidden lg:block hidden" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path d="M12 4v16m8-8H4" />
        </svg>
      </button>

      {/* Spacer */}
      <div className="flex-1" />

      {/* User at bottom */}
      <button className="flex items-center justify-center xl:justify-start gap-3 w-[50px] h-[50px] xl:w-full xl:h-auto xl:p-3 mb-8 rounded-full hover:bg-[var(--qube-surface-hover)] transition-colors">
        <div className="w-10 h-10 rounded-full bg-[var(--qube-surface)] flex items-center justify-center shrink-0">
          <span className="text-sm font-bold text-[var(--qube-text-secondary)]">K</span>
        </div>
        <div className="flex-1 text-left xl:block lg:hidden min-w-0">
          <div className="text-sm font-bold truncate">Kagura</div>
          <div className="text-xs text-[var(--qube-text-secondary)] truncate">@kagura</div>
        </div>
      </button>
    </aside>
  );
}

function NavIcon({ name, active }: { name: string; active: boolean }) {
  const sw = active ? "2.5" : "1.5";
  const cls = "w-[26px] h-[26px] shrink-0";

  switch (name) {
    case "home":
      return <svg className={cls} fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={sw} viewBox="0 0 24 24"><path d="M2.25 12l8.954-8.955a1.126 1.126 0 011.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" /></svg>;
    case "search":
      return <svg className={cls} fill="none" stroke="currentColor" strokeWidth={sw} viewBox="0 0 24 24"><path d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" /></svg>;
    case "bell":
      return <svg className={cls} fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={sw} viewBox="0 0 24 24"><path d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" /></svg>;
    case "mail":
      return <svg className={cls} fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={sw} viewBox="0 0 24 24"><path d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75" /></svg>;
    case "bookmark":
      return <svg className={cls} fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={sw} viewBox="0 0 24 24"><path d="M17.593 3.322c1.1.128 1.907 1.077 1.907 2.185V21L12 17.25 4.5 21V5.507c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0111.186 0z" /></svg>;
    case "user":
      return <svg className={cls} fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth={sw} viewBox="0 0 24 24"><path d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" /></svg>;
    default:
      return null;
  }
}
