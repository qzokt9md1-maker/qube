export default function Home() {
  return (
    <main className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <h1 className="text-6xl font-bold tracking-tight">
          <span className="text-[var(--qube-accent)]">Q</span>ube
        </h1>
        <p className="mt-4 text-lg text-[var(--qube-text-secondary)]">
          A social network where you never miss a post.
        </p>
        <div className="mt-8 flex gap-4 justify-center">
          <button className="rounded-full bg-[var(--qube-primary)] px-8 py-3 font-semibold text-white hover:bg-[var(--qube-primary-dark)] transition-colors">
            Sign Up
          </button>
          <button className="rounded-full border border-[var(--qube-border)] px-8 py-3 font-semibold text-white hover:bg-[var(--qube-surface-hover)] transition-colors">
            Log In
          </button>
        </div>
      </div>
    </main>
  );
}
