export interface User {
  id: string;
  username: string;
  displayName: string;
  bio: string;
  avatarUrl: string;
  headerUrl: string;
  location: string;
  website: string;
  isVerified: boolean;
  isPrivate: boolean;
  followerCount: number;
  followingCount: number;
  postCount: number;
  createdAt: string;
}

export interface Post {
  id: string;
  user: User;
  content: string;
  media: Media[];
  replyToId?: string;
  repostOfId?: string;
  quoteOfId?: string;
  likeCount: number;
  repostCount: number;
  replyCount: number;
  quoteCount: number;
  isLiked: boolean;
  isReposted: boolean;
  isBookmarked: boolean;
  createdAt: string;
}

export interface Media {
  id: string;
  mediaType: "IMAGE" | "VIDEO" | "GIF";
  url: string;
  thumbnailUrl: string;
  width?: number;
  height?: number;
}

export interface Conversation {
  id: string;
  isGroup: boolean;
  name: string;
  participants: User[];
  lastMessage?: Message;
  unreadCount: number;
  updatedAt: string;
}

export interface Message {
  id: string;
  conversationId: string;
  sender: User;
  content: string;
  createdAt: string;
}

export interface Notification {
  id: string;
  actor: User;
  type: "like" | "repost" | "follow" | "reply" | "quote" | "mention" | "dm";
  postId?: string;
  isRead: boolean;
  createdAt: string;
}

export interface TimelineResponse {
  timeline: {
    posts: Post[];
    hasMore: boolean;
    cursor?: string;
    unreadCount: number;
  };
}
