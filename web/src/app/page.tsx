import Link from "next/link";

export default function Home() {
  return (
    <main className="flex min-h-screen items-center justify-center">
      <div className="text-center max-w-lg">
        <h1 className="text-7xl font-bold tracking-tight">
          <span className="text-[var(--qube-accent)]">Q</span>ube
        </h1>
        <p className="mt-4 text-xl text-[var(--qube-text-secondary)]">
          A social network where you never miss a post.
        </p>
        <p className="mt-2 text-sm text-[var(--qube-text-secondary)]">
          Chronological timeline. No algorithm. Every post, in order.
        </p>
        <div className="mt-10 flex gap-4 justify-center">
          <Link
            href="/register"
            className="rounded-full bg-[var(--qube-primary)] px-10 py-3.5 font-semibold text-white hover:bg-[var(--qube-primary-dark)] transition-colors text-lg"
          >
            Sign Up
          </Link>
          <Link
            href="/login"
            className="rounded-full border border-[var(--qube-border)] px-10 py-3.5 font-semibold text-white hover:bg-[var(--qube-surface-hover)] transition-colors text-lg"
          >
            Log In
          </Link>
        </div>
      </div>
    </main>
  );
}
