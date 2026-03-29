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
    if (!content.trim()) return;
    await api.query("CreatePost", { input: { content: content.trim() } });
    setContent("");
    setComposing(false);
    loadTimeline();
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

  return (
    <div>
      {/* Header */}
      <div className="sticky top-0 z-10 bg-[var(--qube-bg)]/80 backdrop-blur-md border-b border-[var(--qube-border)] px-4 py-3">
        <div className="flex items-center justify-between">
          <h1 className="text-xl font-bold">Home</h1>
          {unread > 0 && (
            <span className="bg-[var(--qube-primary)] text-white text-xs px-2 py-1 rounded-full">
              {unread} new
            </span>
          )}
        </div>
      </div>

      {/* Compose */}
      <div className="border-b border-[var(--qube-border)] p-4">
        <div className="flex gap-3">
          <div className="w-11 h-11 rounded-full bg-[var(--qube-surface)] shrink-0" />
          <div className="flex-1">
            <textarea
              placeholder="What's happening?"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              onFocus={() => setComposing(true)}
              className="w-full bg-transparent text-lg resize-none outline-none min-h-[60px] placeholder:text-[var(--qube-text-secondary)]"
              maxLength={500}
            />
            {composing && (
              <div className="flex justify-between items-center mt-2 pt-2 border-t border-[var(--qube-border)]">
                <span className="text-sm text-[var(--qube-text-secondary)]">
                  {500 - content.length}
                </span>
                <button
                  onClick={handlePost}
                  disabled={!content.trim()}
                  className="bg-[var(--qube-primary)] hover:bg-[var(--qube-primary-dark)] text-white font-bold px-5 py-2 rounded-full text-sm disabled:opacity-50 transition-colors"
                >
                  Post
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Timeline */}
      {loading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin w-8 h-8 border-2 border-[var(--qube-accent)] border-t-transparent rounded-full" />
        </div>
      ) : posts.length === 0 ? (
        <div className="text-center py-12 text-[var(--qube-text-secondary)]">
          Follow someone to see their posts!
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
              className="w-full py-4 text-[var(--qube-accent)] hover:bg-[var(--qube-surface-hover)] transition-colors"
            >
              Load more
            </button>
          )}
        </>
      )}
    </div>
  );
}
