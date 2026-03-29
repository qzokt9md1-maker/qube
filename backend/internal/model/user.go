package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	DisplayName    string    `json:"display_name" db:"display_name"`
	Email          string    `json:"email" db:"email"`
	PasswordHash   string    `json:"-" db:"password_hash"`
	Bio            string    `json:"bio" db:"bio"`
	AvatarURL      string    `json:"avatar_url" db:"avatar_url"`
	HeaderURL      string    `json:"header_url" db:"header_url"`
	Location       string    `json:"location" db:"location"`
	Website        string    `json:"website" db:"website"`
	IsVerified     bool      `json:"is_verified" db:"is_verified"`
	IsPrivate      bool      `json:"is_private" db:"is_private"`
	FollowerCount  int       `json:"follower_count" db:"follower_count"`
	FollowingCount int       `json:"following_count" db:"following_count"`
	PostCount      int       `json:"post_count" db:"post_count"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
