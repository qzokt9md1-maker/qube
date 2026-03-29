package model

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	ActorID   uuid.UUID  `json:"actor_id" db:"actor_id"`
	Type      string     `json:"type" db:"type"`
	PostID    *uuid.UUID `json:"post_id" db:"post_id"`
	IsRead    bool       `json:"is_read" db:"is_read"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	// Joined fields
	Actor *User `json:"actor,omitempty"`
	Post  *Post `json:"post,omitempty"`
}
