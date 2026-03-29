package model

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Content      string     `json:"content" db:"content"`
	ReplyToID    *uuid.UUID `json:"reply_to_id" db:"reply_to_id"`
	RepostOfID   *uuid.UUID `json:"repost_of_id" db:"repost_of_id"`
	QuoteOfID    *uuid.UUID `json:"quote_of_id" db:"quote_of_id"`
	LikeCount    int        `json:"like_count" db:"like_count"`
	RepostCount  int        `json:"repost_count" db:"repost_count"`
	ReplyCount   int        `json:"reply_count" db:"reply_count"`
	QuoteCount   int        `json:"quote_count" db:"quote_count"`
	IsDeleted    bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	// Joined fields
	User         *User      `json:"user,omitempty"`
	Media        []Media    `json:"media,omitempty"`
}

type Media struct {
	ID           uuid.UUID `json:"id" db:"id"`
	PostID       uuid.UUID `json:"post_id" db:"post_id"`
	MediaType    string    `json:"media_type" db:"media_type"`
	URL          string    `json:"url" db:"url"`
	ThumbnailURL string    `json:"thumbnail_url" db:"thumbnail_url"`
	Width        *int      `json:"width" db:"width"`
	Height       *int      `json:"height" db:"height"`
	SortOrder    int       `json:"sort_order" db:"sort_order"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
