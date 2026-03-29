package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	Update(ctx context.Context, user *model.User) error
	UpdateCounts(ctx context.Context, id uuid.UUID, field string, delta int) error
	Search(ctx context.Context, query string, limit int, cursor string) ([]*model.User, error)
}

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTimeline(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error)
	GetUserPosts(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error)
	GetReplies(ctx context.Context, postID uuid.UUID, limit int, cursor string) ([]*model.Post, error)
}

type FollowRepository interface {
	Create(ctx context.Context, follow *model.Follow) error
	Delete(ctx context.Context, followerID, followingID uuid.UUID) error
	Exists(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.User, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.User, error)
}

type LikeRepository interface {
	Create(ctx context.Context, userID, postID uuid.UUID) error
	Delete(ctx context.Context, userID, postID uuid.UUID) error
	Exists(ctx context.Context, userID, postID uuid.UUID) (bool, error)
	GetUserLikes(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error)
}

type ConversationRepository interface {
	Create(ctx context.Context, conv *model.Conversation, participantIDs []uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Conversation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Conversation, error)
	GetDMBetween(ctx context.Context, userID1, userID2 uuid.UUID) (*model.Conversation, error)
	MarkRead(ctx context.Context, conversationID, userID uuid.UUID) error
}

type MessageRepository interface {
	Create(ctx context.Context, msg *model.Message) error
	GetByConversation(ctx context.Context, conversationID uuid.UUID, limit int, cursor string) ([]*model.Message, error)
}

type NotificationRepository interface {
	Create(ctx context.Context, notif *model.Notification) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Notification, error)
	MarkRead(ctx context.Context, ids []uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	UnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	GetByToken(ctx context.Context, token string) (*model.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type BookmarkRepository interface {
	Create(ctx context.Context, userID, postID uuid.UUID) error
	Delete(ctx context.Context, userID, postID uuid.UUID) error
	Exists(ctx context.Context, userID, postID uuid.UUID) (bool, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error)
}

type BlockRepository interface {
	Block(ctx context.Context, blockerID, blockedID uuid.UUID) error
	Unblock(ctx context.Context, blockerID, blockedID uuid.UUID) error
	IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error)
}

type MuteRepository interface {
	Mute(ctx context.Context, muterID, mutedID uuid.UUID) error
	Unmute(ctx context.Context, muterID, mutedID uuid.UUID) error
	IsMuted(ctx context.Context, muterID, mutedID uuid.UUID) (bool, error)
}

type TimelineCursorRepository interface {
	Get(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error)
	Update(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error
}
