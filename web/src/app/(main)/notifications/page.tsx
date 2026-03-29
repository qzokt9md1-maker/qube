"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Notification } from "@/types";
import { formatDistanceToNow } from "@/lib/utils";

const typeIcons: Record<string, string> = {
  like: "❤️",
  repost: "🔁",
  follow: "👤",
  reply: "💬",
  quote: "📝",
  mention: "@",
  dm: "✉️",
};

const typeText: Record<string, string> = {
  like: "liked your post",
  repost: "reposted your post",
  follow: "followed you",
  reply: "replied to your post",
  quote: "quoted your post",
  mention: "mentioned you",
  dm: "sent you a message",
};

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    loadNotifications();
  }, []);

  async function loadNotifications() {
    try {
      const data = await api.query("Notifications", { limit: 30 });
      setNotifications(data.notifications.notifications);
      setUnreadCount(data.notifications.unreadCount);
    } catch {}
    setLoading(false);
  }

  async function markAllRead() {
    await api.query("MarkAllNotificationsRead");
    setUnreadCount(0);
    loadNotifications();
  }

  return (
    <div>
      <div className="sticky top-0 z-10 bg-[var(--qube-bg)]/80 backdrop-blur-md border-b border-[var(--qube-border)] px-4 py-3 flex items-center justify-between">
        <h1 className="text-xl font-bold">Notifications</h1>
        {unreadCount > 0 && (
          <button onClick={markAllRead} className="text-sm text-[var(--qube-accent)] hover:underline">
            Mark all read
          </button>
        )}
      </div>

      {loading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin w-8 h-8 border-2 border-[var(--qube-accent)] border-t-transparent rounded-full" />
        </div>
      ) : notifications.length === 0 ? (
        <div className="text-center py-12 text-[var(--qube-text-secondary)]">No notifications yet</div>
      ) : (
        notifications.map((n) => (
          <div
            key={n.id}
            className={`flex items-start gap-3 px-4 py-3 border-b border-[var(--qube-border)] hover:bg-[var(--qube-surface-hover)] cursor-pointer ${
              !n.isRead ? "bg-[var(--qube-accent)]/5" : ""
            }`}
          >
            <span className="text-2xl">{typeIcons[n.type] || "🔔"}</span>
            <div>
              <p className="text-sm">
                <span className="font-bold">{n.actor.displayName}</span>{" "}
                <span className="text-[var(--qube-text-secondary)]">{typeText[n.type]}</span>
              </p>
              <p className="text-xs text-[var(--qube-text-secondary)] mt-1">
                {formatDistanceToNow(n.createdAt)}
              </p>
            </div>
          </div>
        ))
      )}
    </div>
  );
}
