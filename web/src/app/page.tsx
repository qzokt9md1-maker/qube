import Link from "next/link";

export default function Home() {
  return (
    <main className="min-h-screen overflow-hidden">
      {/* Hero */}
      <section className="relative flex min-h-screen items-center justify-center px-6">
        {/* Background effects */}
        <div className="absolute inset-0 overflow-hidden">
          <div className="absolute top-1/4 left-1/4 w-[500px] h-[500px] bg-[var(--qube-primary)] rounded-full opacity-[0.04] blur-[120px]" />
          <div className="absolute bottom-1/4 right-1/4 w-[400px] h-[400px] bg-purple-600 rounded-full opacity-[0.04] blur-[100px]" />
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-indigo-500 rounded-full opacity-[0.03] blur-[150px] animate-pulse-glow" />
        </div>

        <div className="relative text-center max-w-2xl animate-slide-up">
          {/* Logo */}
          <div className="mb-8 animate-float">
            <span className="text-8xl md:text-9xl font-black tracking-tighter">
              <span className="bg-gradient-to-r from-[var(--qube-primary)] via-purple-400 to-[var(--qube-primary-light)] bg-clip-text text-transparent animate-gradient">Q</span>
              <span className="text-white">ube</span>
            </span>
          </div>

          {/* Tagline */}
          <h1 className="text-2xl md:text-3xl font-bold text-white/90 mb-4">
            Never miss what matters.
          </h1>
          <p className="text-lg text-[var(--qube-text-secondary)] mb-4 max-w-md mx-auto">
            The social network with a chronological timeline.<br />
            No algorithm. No noise. Every post, in order.
          </p>

          {/* Key differentiator badges */}
          <div className="flex flex-wrap justify-center gap-3 mb-10">
            <span className="glass px-4 py-2 rounded-full text-sm text-white/80 animate-fade-in delay-100 opacity-0">
              📐 Chronological Feed
            </span>
            <span className="glass px-4 py-2 rounded-full text-sm text-white/80 animate-fade-in delay-200 opacity-0">
              🔖 Unread Tracking
            </span>
            <span className="glass px-4 py-2 rounded-full text-sm text-white/80 animate-fade-in delay-300 opacity-0">
              🚫 No Algorithm
            </span>
            <span className="glass px-4 py-2 rounded-full text-sm text-white/80 animate-fade-in delay-400 opacity-0">
              ⚡ Real-time
            </span>
          </div>

          {/* CTA */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center animate-fade-in delay-500 opacity-0">
            <Link
              href="/register"
              className="rounded-full bg-gradient-to-r from-[var(--qube-primary)] to-purple-500 px-10 py-4 font-bold text-white text-lg hover:shadow-[0_0_30px_rgba(99,102,241,0.4)] hover:scale-105 transition-all duration-200"
            >
              Get Started
            </Link>
            <Link
              href="/login"
              className="rounded-full glass px-10 py-4 font-bold text-white text-lg hover:bg-white/10 transition-all duration-200"
            >
              Log In
            </Link>
          </div>
        </div>

        {/* Scroll indicator */}
        <div className="absolute bottom-8 left-1/2 -translate-x-1/2 animate-bounce">
          <svg className="w-6 h-6 text-[var(--qube-text-secondary)]" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path d="M19 14l-7 7m0 0l-7-7m7 7V3" />
          </svg>
        </div>
      </section>

      {/* Features */}
      <section className="relative px-6 py-24 max-w-5xl mx-auto">
        <h2 className="text-3xl md:text-4xl font-bold text-center mb-4">
          Why <span className="text-[var(--qube-primary)]">Qube</span>?
        </h2>
        <p className="text-center text-[var(--qube-text-secondary)] mb-16 max-w-lg mx-auto">
          Built for people who are tired of missing posts from accounts they actually follow.
        </p>

        <div className="grid md:grid-cols-3 gap-6">
          <FeatureCard
            icon="⏱"
            title="True Chronological"
            description="Posts appear in the exact order they were created. No algorithmic reshuffling, no &quot;suggested&quot; content in your feed."
          />
          <FeatureCard
            icon="📍"
            title="Unread Tracking"
            description="Qube remembers where you left off. Come back anytime and pick up exactly where you stopped — zero posts missed."
          />
          <FeatureCard
            icon="⚡"
            title="Instant Delivery"
            description="Real-time WebSocket push. New posts, DMs, and notifications arrive the moment they happen. No refresh needed."
          />
          <FeatureCard
            icon="🔒"
            title="Your Feed, Your Rules"
            description="No ads, no promoted posts, no engagement bait. Follow who you want, see what they post. That's it."
          />
          <FeatureCard
            icon="💬"
            title="Seamless DMs"
            description="Real-time encrypted messaging. Typing indicators, read receipts, group conversations — all built in."
          />
          <FeatureCard
            icon="🌏"
            title="Open & Global"
            description="Built to scale worldwide. Multi-language support from day one. Your voice deserves a global platform."
          />
        </div>
      </section>

      {/* Comparison */}
      <section className="relative px-6 py-24 border-t border-[var(--qube-border)]">
        <div className="max-w-3xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-12">
            Qube vs The Rest
          </h2>
          <div className="glass rounded-2xl overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white/10">
                  <th className="text-left p-4 text-[var(--qube-text-secondary)]">Feature</th>
                  <th className="p-4 text-[var(--qube-primary)] font-bold">Qube</th>
                  <th className="p-4 text-[var(--qube-text-secondary)]">X</th>
                  <th className="p-4 text-[var(--qube-text-secondary)]">Threads</th>
                </tr>
              </thead>
              <tbody>
                <ComparisonRow feature="Chronological feed" qube="✅ Default" x="🔀 Optional" threads="❌" />
                <ComparisonRow feature="Unread tracking" qube="✅" x="❌" threads="❌" />
                <ComparisonRow feature="No algorithm" qube="✅" x="❌" threads="❌" />
                <ComparisonRow feature="No ads" qube="✅" x="❌" threads="❌" />
                <ComparisonRow feature="Real-time push" qube="✅" x="⚠️ Partial" threads="⚠️ Partial" />
                <ComparisonRow feature="Open API" qube="✅ GraphQL" x="💰 Paid" threads="❌" />
              </tbody>
            </table>
          </div>
        </div>
      </section>

      {/* Final CTA */}
      <section className="relative px-6 py-32 text-center">
        <div className="absolute inset-0 overflow-hidden">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-[var(--qube-primary)] rounded-full opacity-[0.05] blur-[120px]" />
        </div>
        <div className="relative">
          <h2 className="text-4xl md:text-5xl font-bold mb-6">
            Ready to take back<br />your timeline?
          </h2>
          <Link
            href="/register"
            className="inline-block rounded-full bg-gradient-to-r from-[var(--qube-primary)] to-purple-500 px-12 py-4 font-bold text-white text-lg hover:shadow-[0_0_30px_rgba(99,102,241,0.4)] hover:scale-105 transition-all duration-200"
          >
            Join Qube — It&apos;s Free
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-[var(--qube-border)] px-6 py-8 text-center text-sm text-[var(--qube-text-secondary)]">
        <span className="font-bold text-white">Qube</span> &copy; 2026. Built with conviction.
      </footer>
    </main>
  );
}

function FeatureCard({ icon, title, description }: { icon: string; title: string; description: string }) {
  return (
    <div className="glass rounded-2xl p-6 hover:bg-white/[0.04] transition-all duration-200 hover:-translate-y-1 group">
      <div className="text-3xl mb-4">{icon}</div>
      <h3 className="text-lg font-bold mb-2 group-hover:text-[var(--qube-primary)] transition-colors">{title}</h3>
      <p className="text-sm text-[var(--qube-text-secondary)] leading-relaxed">{description}</p>
    </div>
  );
}

function ComparisonRow({ feature, qube, x, threads }: { feature: string; qube: string; x: string; threads: string }) {
  return (
    <tr className="border-b border-white/5 hover:bg-white/[0.02]">
      <td className="p-4 font-medium">{feature}</td>
      <td className="p-4 text-center">{qube}</td>
      <td className="p-4 text-center">{x}</td>
      <td className="p-4 text-center">{threads}</td>
    </tr>
  );
}
