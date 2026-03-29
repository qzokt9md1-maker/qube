"use client";

import { Post } from "@/types";
import { formatDistanceToNow } from "@/lib/utils";
import { useState } from "react";

interface PostCardProps {
  post: Post;
  onLike?: () => void;
  onRepost?: () => void;
  onReply?: () => void;
  onBookmark?: () => void;
  onClick?: () => void;
  onUserClick?: () => void;
}

export function PostCard({ post, onLike, onRepost, onReply, onBookmark, onClick, onUserClick }: PostCardProps) {
  return (
    <article
      className="flex gap-3 px-4 py-3 border-b border-[var(--qube-border)] hover:bg-[var(--qube-surface-hover)] cursor-pointer transition-colors duration-150 animate-fade-in"
      onClick={onClick}
    >
      {/* Avatar */}
      <button
        className="shrink-0 mt-0.5"
        onClick={(e) => { e.stopPropagation(); onUserClick?.(); }}
      >
        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-[var(--qube-primary)] to-purple-600 flex items-center justify-center overflow-hidden ring-0 hover:ring-2 ring-[var(--qube-primary)]/30 transition-all">
          {post.user.avatarUrl ? (
            <img src={post.user.avatarUrl} alt="" className="w-full h-full object-cover" />
          ) : (
            <span className="text-sm font-bold text-white">{post.user.displayName[0]?.toUpperCase()}</span>
          )}
        </div>
      </button>

      {/* Content */}
      <div className="flex-1 min-w-0">
        {/* Header */}
        <div className="flex items-center gap-1 text-[15px] leading-5">
          <span className="font-bold truncate hover:underline">{post.user.displayName}</span>
          {post.user.isVerified && (
            <svg className="w-[18px] h-[18px] text-[var(--qube-accent)] shrink-0" viewBox="0 0 24 24" fill="currentColor">
              <path d="M8.52 3.59a3.04 3.04 0 0 1 6.96 0 3.04 3.04 0 0 1 4.97 4.04 3.04 3.04 0 0 1 0 6.74 3.04 3.04 0 0 1-4.97 4.04 3.04 3.04 0 0 1-6.96 0 3.04 3.04 0 0 1-4.97-4.04 3.04 3.04 0 0 1 0-6.74A3.04 3.04 0 0 1 8.52 3.6Zm7.38 6.2a.75.75 0 1 0-1.06-1.06l-4.59 4.58-1.96-1.97a.75.75 0 0 0-1.06 1.07l2.5 2.5a.75.75 0 0 0 1.06 0l5.12-5.12Z" />
            </svg>
          )}
          <span className="text-[var(--qube-text-secondary)] truncate">@{post.user.username}</span>
          <span className="text-[var(--qube-text-secondary)]">·</span>
          <time className="text-[var(--qube-text-secondary)] whitespace-nowrap hover:underline text-sm">{formatDistanceToNow(post.createdAt)}</time>
        </div>

        {/* Text */}
        <div className="mt-0.5 text-[15px] leading-[22px] whitespace-pre-wrap break-words">{post.content}</div>

        {/* Media */}
        {post.media.length > 0 && (
          <div className="mt-3 rounded-2xl overflow-hidden border border-[var(--qube-border)]">
            {post.media.length === 1 ? (
              <img src={post.media[0].url} alt="" className="w-full max-h-[340px] object-cover" loading="lazy" />
            ) : (
              <div className={`grid gap-0.5 ${post.media.length === 2 ? "grid-cols-2" : post.media.length === 3 ? "grid-cols-2" : "grid-cols-2"}`}>
                {post.media.slice(0, 4).map((m, i) => (
                  <img
                    key={m.id}
                    src={m.url}
                    alt=""
                    className={`w-full object-cover ${post.media.length === 3 && i === 0 ? "row-span-2 h-full" : "h-[170px]"}`}
                    loading="lazy"
                  />
                ))}
              </div>
            )}
          </div>
        )}

        {/* Actions */}
        <div className="flex justify-between mt-2 max-w-[425px] -ml-2">
          <ActionButton icon="reply" count={post.replyCount} hoverColor="text-[var(--qube-primary)] bg-[var(--qube-primary)]/10" onClick={onReply} />
          <ActionButton icon="repost" count={post.repostCount} active={post.isReposted} activeColor="text-emerald-500" hoverColor="text-emerald-500 bg-emerald-500/10" onClick={onRepost} />
          <ActionButton icon="like" count={post.likeCount} active={post.isLiked} activeColor="text-rose-500" hoverColor="text-rose-500 bg-rose-500/10" onClick={onLike} />
          <ActionButton icon="bookmark" active={post.isBookmarked} activeColor="text-[var(--qube-primary)]" hoverColor="text-[var(--qube-primary)] bg-[var(--qube-primary)]/10" onClick={onBookmark} />
          <ActionButton icon="share" hoverColor="text-[var(--qube-primary)] bg-[var(--qube-primary)]/10" />
        </div>
      </div>
    </article>
  );
}

function ActionButton({ icon, count, active, activeColor, hoverColor, onClick }: {
  icon: string;
  count?: number;
  active?: boolean;
  activeColor?: string;
  hoverColor?: string;
  onClick?: () => void;
}) {
  const [pressed, setPressed] = useState(false);

  const baseColor = active ? activeColor : "text-[var(--qube-text-secondary)]";

  return (
    <button
      className={`flex items-center gap-1 text-[13px] p-2 rounded-full transition-all duration-150 group ${baseColor} hover:${hoverColor} ${pressed ? "scale-125" : ""}`}
      onClick={(e) => {
        e.stopPropagation();
        setPressed(true);
        setTimeout(() => setPressed(false), 200);
        onClick?.();
      }}
    >
      <span className={`w-[34px] h-[34px] flex items-center justify-center rounded-full transition-colors group-hover:${hoverColor}`}>
        {icon === "reply" && (
          <svg className="w-[18px] h-[18px]" fill="none" stroke="currentColor" strokeWidth="1.8" viewBox="0 0 24 24">
            <path d="M12 21a9 9 0 1 0-9-9c0 1.5.4 3 1 4.2L3 21l4.8-1c1.2.6 2.7 1 4.2 1Z" />
          </svg>
        )}
        {icon === "repost" && (
          <svg className="w-[18px] h-[18px]" fill="none" stroke="currentColor" strokeWidth="1.8" viewBox="0 0 24 24">
            <path d="M17 2l4 4-4 4M3 11V9a4 4 0 0 1 4-4h14M7 22l-4-4 4-4m14 1v2a4 4 0 0 1-4 4H3" />
          </svg>
        )}
        {icon === "like" && (
          <svg className="w-[18px] h-[18px]" fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth="1.8" viewBox="0 0 24 24">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78Z" />
          </svg>
        )}
        {icon === "bookmark" && (
          <svg className="w-[18px] h-[18px]" fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth="1.8" viewBox="0 0 24 24">
            <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z" />
          </svg>
        )}
        {icon === "share" && (
          <svg className="w-[18px] h-[18px]" fill="none" stroke="currentColor" strokeWidth="1.8" viewBox="0 0 24 24">
            <path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8M16 6l-4-4-4 4M12 2v13" />
          </svg>
        )}
      </span>
      {count !== undefined && count > 0 && <span className="min-w-[1ch]">{formatCount(count)}</span>}
    </button>
  );
}

function formatCount(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
  if (n >= 1000) return `${(n / 1000).toFixed(1)}K`;
  return n.toString();
}
