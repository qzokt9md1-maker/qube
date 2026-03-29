import Link from "next/link";

export default function Home() {
  return (
    <main className="min-h-screen flex">
      {/* Left: branding */}
      <div className="hidden lg:flex flex-1 items-center justify-center">
        <span className="text-[200px] font-black text-[var(--qube-text)] select-none opacity-10">
          Q
        </span>
      </div>

      {/* Right: auth */}
      <div className="flex flex-1 items-center justify-center px-8">
        <div className="max-w-[380px] w-full">
          <h1 className="text-[40px] font-black tracking-tight mb-12 leading-tight">
            Qube
          </h1>
          <h2 className="text-[28px] font-bold mb-2">See everything.</h2>
          <h2 className="text-[28px] font-bold mb-8">Miss nothing.</h2>
          <p className="text-[var(--qube-text-secondary)] mb-10 text-[15px] leading-relaxed">
            Chronological timeline. No algorithm decides what you see. Every post from people you follow, in order.
          </p>

          <div className="flex flex-col gap-3">
            <Link
              href="/register"
              className="w-full bg-[var(--qube-text)] text-black font-bold py-3 rounded-full text-center text-[15px] hover:opacity-90 transition-opacity"
            >
              Create account
            </Link>
            <div className="relative my-1">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-[var(--qube-border)]" />
              </div>
              <div className="relative flex justify-center">
                <span className="bg-[var(--qube-bg)] px-4 text-sm text-[var(--qube-text-secondary)]">or</span>
              </div>
            </div>
            <Link
              href="/login"
              className="w-full border border-[var(--qube-border)] text-[var(--qube-primary)] font-bold py-3 rounded-full text-center text-[15px] hover:bg-[var(--qube-primary)]/10 transition-colors"
            >
              Sign in
            </Link>
          </div>

          <p className="mt-8 text-xs text-[var(--qube-text-secondary)] leading-relaxed">
            By signing up, you agree to the Terms of Service and Privacy Policy.
          </p>
        </div>
      </div>
    </main>
  );
}
