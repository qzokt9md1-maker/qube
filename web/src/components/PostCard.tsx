"use client";

import { Post } from "@/types";
import { formatDistanceToNow } from "@/lib/utils";

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
    <div
      className="flex gap-3 px-4 py-3 border-b border-[var(--qube-border)] hover:bg-[var(--qube-surface-hover)] cursor-pointer transition-colors"
      onClick={onClick}
    >
      {/* Avatar */}
      <div className="shrink-0" onClick={(e) => { e.stopPropagation(); onUserClick?.(); }}>
        <div className="w-11 h-11 rounded-full bg-[var(--qube-surface)] flex items-center justify-center overflow-hidden">
          {post.user.avatarUrl ? (
            <img src={post.user.avatarUrl} alt="" className="w-full h-full object-cover" />
          ) : (
            <span className="text-lg font-bold">{post.user.displayName[0]?.toUpperCase()}</span>
          )}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        {/* Header */}
        <div className="flex items-center gap-1 text-sm">
          <span className="font-bold truncate">{post.user.displayName}</span>
          {post.user.isVerified && (
            <svg className="w-4 h-4 text-[var(--qube-accent)]" fill="currentColor" viewBox="0 0 24 24">
              <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z" />
            </svg>
          )}
          <span className="text-[var(--qube-text-secondary)] truncate">@{post.user.username}</span>
          <span className="text-[var(--qube-text-secondary)]">·</span>
          <span className="text-[var(--qube-text-secondary)] whitespace-nowrap">{formatDistanceToNow(post.createdAt)}</span>
        </div>

        {/* Text */}
        <p className="mt-1 text-[15px] leading-relaxed whitespace-pre-wrap break-words">{post.content}</p>

        {/* Media */}
        {post.media.length > 0 && (
          <div className="mt-2 rounded-2xl overflow-hidden border border-[var(--qube-border)]">
            {post.media.length === 1 ? (
              <img src={post.media[0].url} alt="" className="w-full max-h-[300px] object-cover" />
            ) : (
              <div className="grid grid-cols-2 gap-0.5">
                {post.media.slice(0, 4).map((m) => (
                  <img key={m.id} src={m.url} alt="" className="w-full h-[150px] object-cover" />
                ))}
              </div>
            )}
          </div>
        )}

        {/* Actions */}
        <div className="flex justify-between mt-2 max-w-[400px]">
          <ActionButton icon="reply" count={post.replyCount} onClick={onReply} />
          <ActionButton icon="repost" count={post.repostCount} active={post.isReposted} color="green" onClick={onRepost} />
          <ActionButton icon="like" count={post.likeCount} active={post.isLiked} color="red" onClick={onLike} />
          <ActionButton icon="bookmark" active={post.isBookmarked} color="blue" onClick={onBookmark} />
          <ActionButton icon="share" />
        </div>
      </div>
    </div>
  );
}

function ActionButton({ icon, count, active, color, onClick }: {
  icon: string;
  count?: number;
  active?: boolean;
  color?: string;
  onClick?: () => void;
}) {
  const icons: Record<string, string> = {
    reply: "M1.751 10c0-4.42 3.584-8 8.005-8h4.366c4.49 0 8.129 3.64 8.129 8.13 0 2.25-.893 4.41-2.482 6l-4.127 4.127-.707-.707 4.055-4.055c1.46-1.46 2.26-3.44 2.26-5.37 0-3.94-3.09-7.13-7.13-7.13H9.756c-3.87 0-7.005 3.13-7.005 7v.003L2.75 10zm0 0h.003",
    repost: "M4.5 3.88l4.432 4.14-1.364 1.46L5.5 7.55V16c0 1.1.896 2 2 2H13v2H7.5c-2.209 0-4-1.79-4-4V7.55L1.432 9.48.068 8.02 4.5 3.88zM16.5 6H11V4h5.5c2.209 0 4 1.79 4 4v8.45l2.068-1.93 1.364 1.46-4.432 4.14-4.432-4.14 1.364-1.46 2.068 1.93V8c0-1.1-.896-2-2-2z",
    like: active ? "M20.884 13.19c-1.351 2.48-4.001 5.12-8.379 7.67l-.503.3-.504-.3c-4.379-2.55-7.029-5.19-8.382-7.67-1.36-2.5-1.45-4.92-.334-6.98C3.907 3.85 5.907 2.5 8.245 2.5c1.897 0 3.401.81 4.26 1.63.855-.82 2.357-1.63 4.254-1.63 2.34 0 4.34 1.35 5.47 3.71 1.12 2.06 1.025 4.48-.345 6.98z" : "M16.697 5.5c-1.222-.06-2.679.51-3.89 2.16l-.805 1.09-.806-1.09C9.984 6.01 8.526 5.44 7.304 5.5c-1.243.07-2.349.78-2.91 1.91-.552 1.12-.633 2.78.479 4.82 1.074 1.97 3.257 4.27 7.129 6.61 3.87-2.34 6.052-4.64 7.126-6.61 1.111-2.04 1.03-3.7.477-4.82-.56-1.13-1.666-1.84-2.908-1.91z",
    bookmark: active ? "M4 4.5C4 3.12 5.119 2 6.5 2h11C18.881 2 20 3.12 20 4.5v18.44l-8-5.71-8 5.71V4.5z" : "M4 4.5C4 3.12 5.119 2 6.5 2h11C18.881 2 20 3.12 20 4.5v18.44l-8-5.71-8 5.71V4.5zM6.5 4c-.276 0-.5.22-.5.5v14.56l6-4.29 6 4.29V4.5c0-.28-.224-.5-.5-.5h-11z",
    share: "M12 2.59l5.7 5.7-1.41 1.42L13 6.41V16h-2V6.41l-3.3 3.3-1.41-1.42L12 2.59zM21 15l-.02 3.51c0 1.38-1.12 2.49-2.5 2.49H5.5C4.11 21 3 19.88 3 18.5V15h2v3.5c0 .28.22.5.5.5h12.98c.28 0 .5-.22.5-.5L19 15h2z",
  };

  const activeColors: Record<string, string> = {
    red: "text-red-500",
    green: "text-green-500",
    blue: "text-[var(--qube-accent)]",
  };

  return (
    <button
      className={`flex items-center gap-1 text-xs group ${active && color ? activeColors[color] : "text-[var(--qube-text-secondary)]"} hover:text-[var(--qube-accent)]`}
      onClick={(e) => { e.stopPropagation(); onClick?.(); }}
    >
      <svg className="w-[18px] h-[18px]" fill={active ? "currentColor" : "none"} stroke={active ? "none" : "currentColor"} strokeWidth="1.5" viewBox="0 0 24 24">
        <path d={icons[icon] || ""} />
      </svg>
      {count !== undefined && count > 0 && <span>{formatCount(count)}</span>}
    </button>
  );
}

function formatCount(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
  if (n >= 1000) return `${(n / 1000).toFixed(1)}K`;
  return n.toString();
}
