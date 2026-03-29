"use client";

import { useState, useEffect, useRef } from "react";
import { api } from "@/lib/api";
import { User } from "@/types";

export default function SearchPage() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout>>(undefined);

  useEffect(() => {
    if (!query.trim()) {
      setResults([]);
      return;
    }
    clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => search(query.trim()), 300);
  }, [query]);

  async function search(q: string) {
    setLoading(true);
    try {
      const data = await api.query("SearchUsers", { query: q, limit: 20 });
      setResults(data.searchUsers.users);
    } catch {}
    setLoading(false);
  }

  return (
    <div>
      <div className="sticky top-0 z-10 bg-[var(--qube-bg)]/80 backdrop-blur-md border-b border-[var(--qube-border)] p-3">
        <input
          type="text"
          placeholder="Search users..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="w-full bg-[var(--qube-surface)] border border-[var(--qube-border)] rounded-full py-3 px-5 focus:outline-none focus:border-[var(--qube-accent)]"
          autoFocus
        />
      </div>

      {loading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin w-8 h-8 border-2 border-[var(--qube-accent)] border-t-transparent rounded-full" />
        </div>
      ) : results.length === 0 ? (
        <div className="text-center py-12 text-[var(--qube-text-secondary)]">
          {query ? "No results found" : "Search for users"}
        </div>
      ) : (
        results.map((user) => (
          <div
            key={user.id}
            className="flex items-center gap-3 px-4 py-3 border-b border-[var(--qube-border)] hover:bg-[var(--qube-surface-hover)] cursor-pointer"
          >
            <div className="w-12 h-12 rounded-full bg-[var(--qube-surface)] flex items-center justify-center">
              {user.avatarUrl ? (
                <img src={user.avatarUrl} alt="" className="w-full h-full rounded-full object-cover" />
              ) : (
                <span className="text-lg font-bold">{user.displayName[0]?.toUpperCase()}</span>
              )}
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-1">
                <span className="font-bold">{user.displayName}</span>
                {user.isVerified && <span className="text-[var(--qube-accent)]">✓</span>}
              </div>
              <p className="text-sm text-[var(--qube-text-secondary)]">@{user.username}</p>
              {user.bio && <p className="text-sm mt-1 text-[var(--qube-text-secondary)] truncate">{user.bio}</p>}
            </div>
            <span className="text-xs text-[var(--qube-text-secondary)]">{user.followerCount} followers</span>
          </div>
        ))
      )}
    </div>
  );
}
