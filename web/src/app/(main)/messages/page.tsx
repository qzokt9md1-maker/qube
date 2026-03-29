"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Conversation } from "@/types";
import { formatDistanceToNow } from "@/lib/utils";

export default function MessagesPage() {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadConversations();
  }, []);

  async function loadConversations() {
    try {
      const data = await api.query("Conversations", { limit: 20 });
      setConversations(data.conversations.conversations);
    } catch {}
    setLoading(false);
  }

  return (
    <div>
      <div className="sticky top-0 z-10 bg-[var(--qube-bg)]/80 backdrop-blur-md border-b border-[var(--qube-border)] px-4 py-3 flex items-center justify-between">
        <h1 className="text-xl font-bold">Messages</h1>
        <button className="text-[var(--qube-accent)] hover:bg-[var(--qube-surface-hover)] p-2 rounded-full">
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>

      {loading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin w-8 h-8 border-2 border-[var(--qube-accent)] border-t-transparent rounded-full" />
        </div>
      ) : conversations.length === 0 ? (
        <div className="text-center py-12 text-[var(--qube-text-secondary)]">
          No messages yet. Start a conversation!
        </div>
      ) : (
        conversations.map((conv) => {
          const other = conv.participants[0];
          return (
            <div
              key={conv.id}
              className="flex items-center gap-3 px-4 py-3 border-b border-[var(--qube-border)] hover:bg-[var(--qube-surface-hover)] cursor-pointer"
            >
              <div className="w-12 h-12 rounded-full bg-[var(--qube-surface)] flex items-center justify-center shrink-0">
                {other?.avatarUrl ? (
                  <img src={other.avatarUrl} alt="" className="w-full h-full rounded-full object-cover" />
                ) : (
                  <span className="text-lg font-bold">{other?.displayName?.[0]?.toUpperCase() || "?"}</span>
                )}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between">
                  <span className="font-bold truncate">{conv.isGroup ? conv.name : other?.displayName}</span>
                  <span className="text-xs text-[var(--qube-text-secondary)]">{formatDistanceToNow(conv.updatedAt)}</span>
                </div>
                {conv.lastMessage && (
                  <p className="text-sm text-[var(--qube-text-secondary)] truncate">{conv.lastMessage.content}</p>
                )}
              </div>
              {conv.unreadCount > 0 && (
                <span className="bg-[var(--qube-primary)] text-white text-xs w-5 h-5 rounded-full flex items-center justify-center">
                  {conv.unreadCount}
                </span>
              )}
            </div>
          );
        })
      )}
    </div>
  );
}
