"use client";

import { useState } from "react";

export function RightSidebar() {
  const [searchFocused, setSearchFocused] = useState(false);

  return (
    <aside className="fixed right-0 top-0 h-full w-[350px] p-4 overflow-y-auto max-xl:hidden">
      {/* Search */}
      <div className="sticky top-4">
        <div className={`relative transition-all duration-200 ${searchFocused ? "ring-1 ring-[var(--qube-primary)]" : ""} rounded-full`}>
          <input
            type="text"
            placeholder="Search Qube"
            className="w-full bg-[var(--qube-surface)] border border-[var(--qube-border)] rounded-full py-3 px-12 text-sm focus:outline-none focus:border-[var(--qube-primary)] focus:bg-transparent transition-colors"
            onFocus={() => setSearchFocused(true)}
            onBlur={() => setSearchFocused(false)}
          />
          <svg className={`absolute left-4 top-3.5 w-5 h-5 transition-colors ${searchFocused ? "text-[var(--qube-primary)]" : "text-[var(--qube-text-secondary)]"}`} fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>

        {/* Trending */}
        <div className="mt-4 bg-[var(--qube-surface)] rounded-2xl overflow-hidden">
          <h2 className="text-xl font-bold p-4 pb-2">Trending</h2>
          {[
            { tag: "#Qube", category: "Technology", posts: "1,234" },
            { tag: "#tech", category: "Technology", posts: "892" },
            { tag: "#coding", category: "Development", posts: "567" },
            { tag: "#startup", category: "Business", posts: "345" },
            { tag: "#design", category: "Design", posts: "234" },
          ].map((item) => (
            <div
              key={item.tag}
              className="px-4 py-3 hover:bg-[var(--qube-surface-hover)] cursor-pointer transition-colors"
            >
              <p className="text-xs text-[var(--qube-text-secondary)]">{item.category} · Trending</p>
              <p className="font-bold text-[15px] mt-0.5">{item.tag}</p>
              <p className="text-xs text-[var(--qube-text-secondary)] mt-0.5">{item.posts} posts</p>
            </div>
          ))}
          <button className="w-full text-left px-4 py-3 text-[var(--qube-primary)] text-sm hover:bg-[var(--qube-surface-hover)] transition-colors">
            Show more
          </button>
        </div>

        {/* Who to follow */}
        <div className="mt-4 bg-[var(--qube-surface)] rounded-2xl overflow-hidden">
          <h2 className="text-xl font-bold p-4 pb-2">Who to follow</h2>
          {[
            { name: "Qube Official", handle: "@qube", initial: "Q" },
            { name: "Tech News", handle: "@technews", initial: "T" },
            { name: "Design Daily", handle: "@designdaily", initial: "D" },
          ].map((user) => (
            <div
              key={user.handle}
              className="flex items-center gap-3 px-4 py-3 hover:bg-[var(--qube-surface-hover)] cursor-pointer transition-colors"
            >
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[var(--qube-primary)] to-purple-600 flex items-center justify-center shrink-0">
                <span className="text-sm font-bold">{user.initial}</span>
              </div>
              <div className="flex-1 min-w-0">
                <p className="font-bold text-sm truncate">{user.name}</p>
                <p className="text-xs text-[var(--qube-text-secondary)]">{user.handle}</p>
              </div>
              <button className="bg-white text-black font-bold text-sm px-4 py-1.5 rounded-full hover:bg-white/90 transition-colors">
                Follow
              </button>
            </div>
          ))}
          <button className="w-full text-left px-4 py-3 text-[var(--qube-primary)] text-sm hover:bg-[var(--qube-surface-hover)] transition-colors">
            Show more
          </button>
        </div>

        {/* Footer */}
        <div className="mt-4 px-4 text-xs text-[var(--qube-text-secondary)] flex flex-wrap gap-x-3 gap-y-1">
          <a href="#" className="hover:underline">Terms</a>
          <a href="#" className="hover:underline">Privacy</a>
          <a href="#" className="hover:underline">About</a>
          <a href="#" className="hover:underline">Accessibility</a>
          <span>&copy; 2026 Qube</span>
        </div>
      </div>
    </aside>
  );
}
