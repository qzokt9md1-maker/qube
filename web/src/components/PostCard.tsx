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
    <article
      className="flex gap-3 px-4 py-3 border-b border-[var(--qube-border)] hover:bg-white/[0.03] cursor-pointer transition-colors"
      onClick={onClick}
    >
      {/* Avatar */}
      <button
        className="shrink-0 mt-0.5"
        onClick={(e) => { e.stopPropagation(); onUserClick?.(); }}
      >
        <div className="w-10 h-10 rounded-full bg-[var(--qube-surface)] flex items-center justify-center overflow-hidden">
          {post.user.avatarUrl ? (
            <img src={post.user.avatarUrl} alt="" className="w-full h-full object-cover" />
          ) : (
            <span className="text-sm font-bold text-[var(--qube-text-secondary)]">{post.user.displayName[0]?.toUpperCase()}</span>
          )}
        </div>
      </button>

      <div className="flex-1 min-w-0">
        {/* Header */}
        <div className="flex items-center gap-1 text-[15px] leading-5">
          <span className="font-bold truncate">{post.user.displayName}</span>
          {post.user.isVerified && (
            <svg className="w-[18px] h-[18px] text-[var(--qube-primary)] shrink-0" viewBox="0 0 24 24" fill="currentColor">
              <path d="M8.52 3.59a3.04 3.04 0 0 1 6.96 0 3.04 3.04 0 0 1 4.97 4.04 3.04 3.04 0 0 1 0 6.74 3.04 3.04 0 0 1-4.97 4.04 3.04 3.04 0 0 1-6.96 0 3.04 3.04 0 0 1-4.97-4.04 3.04 3.04 0 0 1 0-6.74A3.04 3.04 0 0 1 8.52 3.6Zm7.38 6.2a.75.75 0 1 0-1.06-1.06l-4.59 4.58-1.96-1.97a.75.75 0 0 0-1.06 1.07l2.5 2.5a.75.75 0 0 0 1.06 0l5.12-5.12Z" />
            </svg>
          )}
          <span className="text-[var(--qube-text-secondary)] truncate font-normal">@{post.user.username}</span>
          <span className="text-[var(--qube-text-secondary)]">·</span>
          <time className="text-[var(--qube-text-secondary)] whitespace-nowrap text-sm font-normal">{formatDistanceToNow(post.createdAt)}</time>
        </div>

        {/* Content - the star of the show */}
        <div className="mt-0.5 text-[15px] leading-[20px] whitespace-pre-wrap break-words">{post.content}</div>

        {/* Media */}
        {post.media.length > 0 && (
          <div className="mt-3 rounded-2xl overflow-hidden border border-[var(--qube-border)]">
            {post.media.length === 1 ? (
              <img src={post.media[0].url} alt="" className="w-full max-h-[340px] object-cover" loading="lazy" />
            ) : (
              <div className="grid grid-cols-2 gap-0.5">
                {post.media.slice(0, 4).map((m) => (
                  <img key={m.id} src={m.url} alt="" className="w-full h-[170px] object-cover" loading="lazy" />
                ))}
              </div>
            )}
          </div>
        )}

        {/* Actions - all gray by default, colored only when active */}
        <div className="flex justify-between mt-1 max-w-[425px] -ml-2">
          <ActionBtn icon="reply" count={post.replyCount} onClick={onReply} />
          <ActionBtn icon="repost" count={post.repostCount} active={post.isReposted} activeColor="var(--qube-repost)" onClick={onRepost} />
          <ActionBtn icon="like" count={post.likeCount} active={post.isLiked} activeColor="var(--qube-like)" onClick={onLike} />
          <ActionBtn icon="bookmark" active={post.isBookmarked} activeColor="var(--qube-primary)" onClick={onBookmark} />
          <ActionBtn icon="share" onClick={() => {}} />
        </div>
      </div>
    </article>
  );
}

function ActionBtn({ icon, count, active, activeColor, onClick }: {
  icon: string; count?: number; active?: boolean; activeColor?: string; onClick?: () => void;
}) {
  const color = active && activeColor ? activeColor : "var(--qube-text-secondary)";

  return (
    <button
      className="flex items-center gap-1 p-2 rounded-full group transition-colors"
      style={{ color }}
      onClick={(e) => { e.stopPropagation(); onClick?.(); }}
    >
      <span className="w-[34px] h-[34px] flex items-center justify-center rounded-full group-hover:bg-white/[0.08] transition-colors">
        {icon === "reply" && (
          <svg className="w-[18px] h-[18px]" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
            <path d="M12 21a9 9 0 1 0-9-9c0 1.5.4 3 1 4.2L3 21l4.8-1c1.2.6 2.7 1 4.2 1Z" />
          </svg>
        )}
        {icon === "repost" && (
          <svg className="w-[18px] h-[18px]" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
            <path d="M17 2l4 4-4 4M3 11V9a4 4 0 0 1 4-4h14M7 22l-4-4 4-4m14 1v2a4 4 0 0 1-4 4H3" />
          </svg>
        )}
        {icon === "like" && (
          <svg className="w-[18px] h-[18px]" fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78Z" />
          </svg>
        )}
        {icon === "bookmark" && (
          <svg className="w-[18px] h-[18px]" fill={active ? "currentColor" : "none"} stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
            <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z" />
          </svg>
        )}
        {icon === "share" && (
          <svg className="w-[18px] h-[18px]" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
            <path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8M16 6l-4-4-4 4M12 2v13" />
          </svg>
        )}
      </span>
      {count !== undefined && count > 0 && (
        <span className="text-[13px]">{count >= 1000 ? `${(count / 1000).toFixed(1)}K` : count}</span>
      )}
    </button>
  );
}
