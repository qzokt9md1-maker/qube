"use client";

export function RightSidebar() {
  return (
    <aside className="fixed right-0 top-0 h-full w-[350px] border-l border-[var(--qube-border)] p-4">
      {/* Search */}
      <div className="relative">
        <input
          type="text"
          placeholder="Search Qube"
          className="w-full bg-[var(--qube-surface)] border border-[var(--qube-border)] rounded-full py-3 px-5 pl-12 text-sm focus:outline-none focus:border-[var(--qube-accent)]"
        />
        <svg className="absolute left-4 top-3.5 w-5 h-5 text-[var(--qube-text-secondary)]" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>

      {/* Trending */}
      <div className="mt-4 bg-[var(--qube-surface)] rounded-2xl p-4">
        <h2 className="text-xl font-bold mb-3">Trending</h2>
        <div className="space-y-4">
          {["#Qube", "#tech", "#coding"].map((tag) => (
            <div key={tag} className="hover:bg-[var(--qube-surface-hover)] -mx-4 px-4 py-2 cursor-pointer">
              <p className="text-[var(--qube-text-secondary)] text-xs">Trending</p>
              <p className="font-bold">{tag}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Footer */}
      <div className="mt-4 text-xs text-[var(--qube-text-secondary)] space-x-2">
        <span>Terms</span>
        <span>Privacy</span>
        <span>About</span>
        <span>© 2026 Qube</span>
      </div>
    </aside>
  );
}
