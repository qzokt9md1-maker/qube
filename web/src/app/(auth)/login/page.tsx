"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { api } from "@/lib/api";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const data = await api.query("Login", {
        input: { email, password },
      });
      api.setTokens(data.login.accessToken, data.login.refreshToken);
      router.push("/home");
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center">
      <div className="w-full max-w-sm p-8">
        <h1 className="text-3xl font-bold text-center mb-8">
          <span className="text-[var(--qube-accent)]">Q</span>ube
        </h1>
        <h2 className="text-xl font-bold mb-6">Log in</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full bg-[var(--qube-surface)] border border-[var(--qube-border)] rounded-lg p-3 focus:outline-none focus:border-[var(--qube-accent)]"
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full bg-[var(--qube-surface)] border border-[var(--qube-border)] rounded-lg p-3 focus:outline-none focus:border-[var(--qube-accent)]"
            required
          />
          {error && <p className="text-red-500 text-sm">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-[var(--qube-primary)] hover:bg-[var(--qube-primary-dark)] text-white font-bold py-3 rounded-full transition-colors disabled:opacity-50"
          >
            {loading ? "Logging in..." : "Log In"}
          </button>
        </form>
        <p className="mt-6 text-center text-[var(--qube-text-secondary)]">
          Don&apos;t have an account?{" "}
          <Link href="/register" className="text-[var(--qube-accent)] hover:underline">
            Sign up
          </Link>
        </p>
      </div>
    </main>
  );
}
