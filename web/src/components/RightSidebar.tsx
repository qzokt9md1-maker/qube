"use client";

export function RightSidebar() {
  return (
    <aside className="fixed right-0 top-0 h-full w-[350px] p-4 overflow-y-auto max-xl:hidden">
      <div className="sticky top-4">
        {/* Search */}
        <div className="relative">
          <input
            type="text"
            placeholder="Search Qube"
            className="w-full bg-[var(--qube-surface)] border border-[var(--qube-border)] rounded-full py-2.5 px-12 text-sm focus:outline-none focus:border-[var(--qube-primary)] transition-colors"
          />
          <svg className="absolute left-4 top-3 w-4 h-4 text-[var(--qube-text-secondary)]" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>

        {/* Trending */}
        <div className="mt-4 bg-[var(--qube-surface)] rounded-2xl overflow-hidden">
          <h2 className="font-bold text-xl px-4 pt-3 pb-2">Trending</h2>
          {[
            { tag: "#Qube", posts: "1,234" },
            { tag: "#tech", posts: "892" },
            { tag: "#coding", posts: "567" },
          ].map((item) => (
            <div key={item.tag} className="px-4 py-2.5 hover:bg-white/[0.03] cursor-pointer transition-colors">
              <p className="font-bold text-[15px]">{item.tag}</p>
              <p className="text-xs text-[var(--qube-text-secondary)]">{item.posts} posts</p>
            </div>
          ))}
          <div className="px-4 py-3 text-[var(--qube-primary)] text-sm cursor-pointer hover:bg-white/[0.03] transition-colors">
            Show more
          </div>
        </div>

        {/* Who to follow */}
        <div className="mt-4 bg-[var(--qube-surface)] rounded-2xl overflow-hidden">
          <h2 className="font-bold text-xl px-4 pt-3 pb-2">Who to follow</h2>
          {[
            { name: "Qube Official", handle: "@qube" },
            { name: "Tech News", handle: "@technews" },
          ].map((user) => (
            <div key={user.handle} className="flex items-center gap-3 px-4 py-2.5 hover:bg-white/[0.03] cursor-pointer transition-colors">
              <div className="w-10 h-10 rounded-full bg-[var(--qube-border)] flex items-center justify-center shrink-0">
                <span className="text-sm text-[var(--qube-text-secondary)]">{user.name[0]}</span>
              </div>
              <div className="flex-1 min-w-0">
                <p className="font-bold text-sm truncate">{user.name}</p>
                <p className="text-xs text-[var(--qube-text-secondary)]">{user.handle}</p>
              </div>
              <button className="bg-[var(--qube-text)] text-black font-bold text-sm px-4 py-1.5 rounded-full hover:opacity-90 transition-opacity">
                Follow
              </button>
            </div>
          ))}
        </div>

        <div className="mt-4 px-4 text-xs text-[var(--qube-text-secondary)] flex flex-wrap gap-x-3 gap-y-1">
          <span>Terms</span><span>Privacy</span><span>About</span><span>&copy; 2026 Qube</span>
        </div>
      </div>
    </aside>
  );
}
