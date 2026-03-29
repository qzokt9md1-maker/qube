"use client";

import { useEffect, useState, useCallback } from "react";
import { api } from "@/lib/api";
import { Post } from "@/types";
import { PostCard } from "@/components/PostCard";

export default function HomePage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [cursor, setCursor] = useState<string>();
  const [hasMore, setHasMore] = useState(true);
  const [unread, setUnread] = useState(0);
  const [composing, setComposing] = useState(false);
  const [content, setContent] = useState("");
  const [posting, setPosting] = useState(false);

  const loadTimeline = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.query("Timeline", { limit: 20 });
      setPosts(data.timeline.posts);
      setHasMore(data.timeline.hasMore);
      setCursor(data.timeline.cursor);
      setUnread(data.timeline.unreadCount);
      if (data.timeline.posts.length > 0) {
        api.query("UpdateTimelineCursor", { lastSeenPostId: data.timeline.posts[0].id });
      }
    } catch {}
    setLoading(false);
  }, []);

  useEffect(() => {
    loadTimeline();
  }, [loadTimeline]);

  async function loadMore() {
    if (!hasMore) return;
    const data = await api.query("Timeline", { limit: 20, cursor });
    setPosts((prev) => [...prev, ...data.timeline.posts]);
    setHasMore(data.timeline.hasMore);
    setCursor(data.timeline.cursor);
  }

  async function handlePost() {
    if (!content.trim() || posting) return;
    setPosting(true);
    try {
      await api.query("CreatePost", { input: { content: content.trim() } });
      setContent("");
      setComposing(false);
      await loadTimeline();
    } catch {}
    setPosting(false);
  }

  async function handleLike(post: Post) {
    const op = post.isLiked ? "UnlikePost" : "LikePost";
    await api.query(op, { postId: post.id });
    loadTimeline();
  }

  async function handleRepost(post: Post) {
    await api.query("Repost", { postId: post.id });
    loadTimeline();
  }

  const charCount = 500 - content.length;

  return (
    <div>
      {/* Header */}
      <div className="sticky top-0 z-10 bg-[var(--qube-bg)]/80 backdrop-blur-xl border-b border-[var(--qube-border)]">
        <div className="flex items-center justify-between px-4 h-[53px]">
          <h1 className="text-xl font-bold">Home</h1>
          {unread > 0 && (
            <button
              onClick={loadTimeline}
              className="bg-[var(--qube-primary)] text-white text-xs font-bold px-3 py-1.5 rounded-full hover:bg-[var(--qube-primary-dark)] transition-colors animate-fade-in"
            >
              {unread} new post{unread > 1 ? "s" : ""}
            </button>
          )}
        </div>
      </div>

      {/* Compose */}
      <div className="border-b border-[var(--qube-border)] px-4 pt-3 pb-2">
        <div className="flex gap-3">
          <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[var(--qube-primary)] to-purple-600 shrink-0 flex items-center justify-center">
            <span className="text-sm font-bold">K</span>
          </div>
          <div className="flex-1">
            <textarea
              placeholder="What's happening?"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              onFocus={() => setComposing(true)}
              className="w-full bg-transparent text-[17px] resize-none outline-none min-h-[52px] placeholder:text-[var(--qube-text-secondary)] leading-relaxed"
              maxLength={500}
              onKeyDown={(e) => {
                if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) handlePost();
              }}
            />
            <div className={`flex justify-between items-center pt-2 border-t border-[var(--qube-border)] transition-all ${composing ? "opacity-100 mb-1" : "opacity-0 h-0 overflow-hidden"}`}>
              <div className="flex items-center gap-2">
                <button className="p-2 rounded-full hover:bg-[var(--qube-primary)]/10 text-[var(--qube-primary)] transition-colors">
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
                    <path d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909M3.75 21h16.5a1.5 1.5 0 001.5-1.5V4.5a1.5 1.5 0 00-1.5-1.5H3.75a1.5 1.5 0 00-1.5 1.5v15a1.5 1.5 0 001.5 1.5z" />
                  </svg>
                </button>
                <button className="p-2 rounded-full hover:bg-[var(--qube-primary)]/10 text-[var(--qube-primary)] transition-colors">
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
                    <path d="M15.182 15.182a4.5 4.5 0 01-6.364 0M21 12a9 9 0 11-18 0 9 9 0 0118 0zM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75zm-.375 0h.008v.015h-.008V9.75zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75zm-.375 0h.008v.015h-.008V9.75z" />
                  </svg>
                </button>
              </div>
              <div className="flex items-center gap-3">
                {content.length > 0 && (
                  <div className="relative w-7 h-7">
                    <svg className="w-7 h-7 -rotate-90" viewBox="0 0 32 32">
                      <circle cx="16" cy="16" r="14" fill="none" stroke="var(--qube-border)" strokeWidth="2" />
                      <circle
                        cx="16" cy="16" r="14" fill="none"
                        stroke={charCount < 0 ? "var(--qube-danger)" : charCount < 20 ? "#f59e0b" : "var(--qube-primary)"}
                        strokeWidth="2"
                        strokeDasharray={`${Math.max(0, (content.length / 500)) * 87.96} 87.96`}
                        strokeLinecap="round"
                      />
                    </svg>
                    {charCount <= 20 && (
                      <span className={`absolute inset-0 flex items-center justify-center text-[10px] font-bold ${charCount < 0 ? "text-[var(--qube-danger)]" : "text-[var(--qube-text-secondary)]"}`}>
                        {charCount}
                      </span>
                    )}
                  </div>
                )}
                <button
                  onClick={handlePost}
                  disabled={!content.trim() || posting || charCount < 0}
                  className="bg-[var(--qube-primary)] hover:bg-[var(--qube-primary-dark)] text-white font-bold px-5 py-1.5 rounded-full text-sm disabled:opacity-40 transition-all hover:shadow-[0_0_15px_rgba(99,102,241,0.3)]"
                >
                  {posting ? (
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  ) : (
                    "Post"
                  )}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Timeline */}
      {loading ? (
        <div className="flex justify-center py-16">
          <div className="w-8 h-8 border-2 border-[var(--qube-primary)] border-t-transparent rounded-full animate-spin" />
        </div>
      ) : posts.length === 0 ? (
        <div className="text-center py-20 px-8">
          <div className="text-5xl mb-4">👋</div>
          <h3 className="text-xl font-bold mb-2">Welcome to Qube!</h3>
          <p className="text-[var(--qube-text-secondary)] max-w-sm mx-auto">
            Follow people you&apos;re interested in to see their posts here. Your timeline is 100% chronological — no algorithm, no surprises.
          </p>
        </div>
      ) : (
        <>
          {posts.map((post) => (
            <PostCard
              key={post.id}
              post={post}
              onLike={() => handleLike(post)}
              onRepost={() => handleRepost(post)}
            />
          ))}
          {hasMore && (
            <button
              onClick={loadMore}
              className="w-full py-4 text-[var(--qube-primary)] hover:bg-[var(--qube-surface-hover)] transition-colors font-medium"
            >
              Show more
            </button>
          )}
        </>
      )}
    </div>
  );
}
