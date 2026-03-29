"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { api } from "@/lib/api";

export default function RegisterPage() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const data = await api.query("Register", {
        input: { username, displayName, email, password },
      });
      api.setTokens(data.register.accessToken, data.register.refreshToken);
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
        <h2 className="text-xl font-bold mb-6">Create your account</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full bg-[var(--qube-surface)] border border-[var(--qube-border)] rounded-lg p-3 focus:outline-none focus:border-[var(--qube-accent)]"
            required
          />
          <input
            type="text"
            placeholder="Display Name"
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            className="w-full bg-[var(--qube-surface)] border border-[var(--qube-border)] rounded-lg p-3 focus:outline-none focus:border-[var(--qube-accent)]"
            required
          />
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
            placeholder="Password (8+ characters)"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            minLength={8}
            className="w-full bg-[var(--qube-surface)] border border-[var(--qube-border)] rounded-lg p-3 focus:outline-none focus:border-[var(--qube-accent)]"
            required
          />
          {error && <p className="text-red-500 text-sm">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-[var(--qube-primary)] hover:bg-[var(--qube-primary-dark)] text-white font-bold py-3 rounded-full transition-colors disabled:opacity-50"
          >
            {loading ? "Creating..." : "Create Account"}
          </button>
        </form>
        <p className="mt-6 text-center text-[var(--qube-text-secondary)]">
          Already have an account?{" "}
          <Link href="/login" className="text-[var(--qube-accent)] hover:underline">
            Log in
          </Link>
        </p>
      </div>
    </main>
  );
}
