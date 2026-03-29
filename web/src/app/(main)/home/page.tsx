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

  useEffect(() => { loadTimeline(); }, [loadTimeline]);

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
      await loadTimeline();
    } catch {}
    setPosting(false);
  }

  async function handleLike(post: Post) {
    await api.query(post.isLiked ? "UnlikePost" : "LikePost", { postId: post.id });
    loadTimeline();
  }

  async function handleRepost(post: Post) {
    await api.query("Repost", { postId: post.id });
    loadTimeline();
  }

  return (
    <div>
      {/* Header */}
      <div className="sticky top-0 z-10 bg-[var(--qube-bg)]/80 backdrop-blur-md border-b border-[var(--qube-border)]">
        <div className="flex items-center justify-between px-4 h-[53px]">
          <h1 className="text-xl font-bold"><span className="text-[var(--qube-primary)]">Q</span>ube</h1>
          {unread > 0 && (
            <button onClick={loadTimeline} className="text-[var(--qube-primary)] text-sm font-medium hover:underline">
              {unread} new post{unread > 1 ? "s" : ""}
            </button>
          )}
        </div>
      </div>

      {/* Compose */}
      <div className="border-b border-[var(--qube-border)] px-4 pt-3 pb-2">
        <div className="flex gap-3">
          <div className="w-10 h-10 rounded-full bg-[var(--qube-surface)] shrink-0 flex items-center justify-center">
            <span className="text-sm text-[var(--qube-text-secondary)]">K</span>
          </div>
          <div className="flex-1">
            <textarea
              placeholder="What's happening?"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              className="w-full bg-transparent text-[17px] resize-none outline-none min-h-[52px] placeholder:text-[var(--qube-text-secondary)]"
              maxLength={500}
              onKeyDown={(e) => { if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) handlePost(); }}
            />
            {content.length > 0 && (
              <div className="flex justify-between items-center pt-2 border-t border-[var(--qube-border)]">
                <span className={`text-sm ${500 - content.length < 20 ? "text-[var(--qube-danger)]" : "text-[var(--qube-text-secondary)]"}`}>
                  {500 - content.length}
                </span>
                <button
                  onClick={handlePost}
                  disabled={!content.trim() || posting}
                  className="bg-[var(--qube-primary)] hover:bg-[var(--qube-primary-dark)] text-white font-bold px-5 py-1.5 rounded-full text-sm disabled:opacity-40 transition-colors"
                >
                  {posting ? "..." : "Post"}
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Timeline */}
      {loading ? (
        <div className="flex justify-center py-16">
          <div className="w-6 h-6 border-2 border-[var(--qube-text-secondary)] border-t-transparent rounded-full animate-spin" />
        </div>
      ) : posts.length === 0 ? (
        <div className="text-center py-20 px-8 text-[var(--qube-text-secondary)]">
          <p className="text-lg mb-1">Your timeline is empty</p>
          <p className="text-sm">Follow people to see their posts here.</p>
        </div>
      ) : (
        <>
          {posts.map((post) => (
            <PostCard key={post.id} post={post} onLike={() => handleLike(post)} onRepost={() => handleRepost(post)} />
          ))}
          {hasMore && (
            <button onClick={loadMore} className="w-full py-4 text-[var(--qube-primary)] hover:bg-white/[0.03] transition-colors text-sm">
              Show more
            </button>
          )}
        </>
      )}
    </div>
  );
}
