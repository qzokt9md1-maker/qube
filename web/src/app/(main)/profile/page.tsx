"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { User, Post } from "@/types";
import { PostCard } from "@/components/PostCard";

export default function ProfilePage() {
  const [user, setUser] = useState<User | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState<"posts" | "replies" | "likes">("posts");

  useEffect(() => {
    loadProfile();
  }, []);

  async function loadProfile() {
    try {
      const data = await api.query("Me");
      setUser(data.me);
      const postsData = await api.query("UserPosts", { username: data.me.username, limit: 20 });
      setPosts(postsData.userPosts.posts);
    } catch {}
    setLoading(false);
  }

  if (loading) {
    return (
      <div className="flex justify-center py-16">
        <div className="animate-spin w-8 h-8 border-2 border-[var(--qube-accent)] border-t-transparent rounded-full" />
      </div>
    );
  }

  if (!user) return <div className="text-center py-12">Not logged in</div>;

  return (
    <div>
      {/* Header */}
      <div className="h-[200px] bg-gradient-to-br from-[var(--qube-primary)] to-purple-900" />

      {/* Profile info */}
      <div className="px-4 pb-4">
        <div className="flex justify-between items-end -mt-16">
          <div className="w-[134px] h-[134px] rounded-full border-4 border-[var(--qube-bg)] bg-[var(--qube-surface)] flex items-center justify-center overflow-hidden">
            {user.avatarUrl ? (
              <img src={user.avatarUrl} alt="" className="w-full h-full object-cover" />
            ) : (
              <span className="text-5xl font-bold">{user.displayName[0]?.toUpperCase()}</span>
            )}
          </div>
          <button className="border border-[var(--qube-border)] rounded-full px-5 py-2 font-bold hover:bg-[var(--qube-surface-hover)] transition-colors">
            Edit profile
          </button>
        </div>

        <div className="mt-3">
          <div className="flex items-center gap-1">
            <h2 className="text-xl font-bold">{user.displayName}</h2>
            {user.isVerified && <span className="text-[var(--qube-accent)]">✓</span>}
          </div>
          <p className="text-[var(--qube-text-secondary)]">@{user.username}</p>
          {user.bio && <p className="mt-3 text-[15px]">{user.bio}</p>}
          <div className="flex gap-5 mt-3 text-sm">
            <span>
              <strong>{user.followingCount}</strong>{" "}
              <span className="text-[var(--qube-text-secondary)]">Following</span>
            </span>
            <span>
              <strong>{user.followerCount}</strong>{" "}
              <span className="text-[var(--qube-text-secondary)]">Followers</span>
            </span>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-[var(--qube-border)]">
        {(["posts", "replies", "likes"] as const).map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`flex-1 py-4 text-center capitalize hover:bg-[var(--qube-surface-hover)] transition-colors ${
              tab === t ? "font-bold border-b-2 border-[var(--qube-accent)]" : "text-[var(--qube-text-secondary)]"
            }`}
          >
            {t}
          </button>
        ))}
      </div>

      {/* Content */}
      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  );
}
