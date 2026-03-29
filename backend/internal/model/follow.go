package model

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uuid.UUID `json:"id" db:"id"`
	FollowerID  uuid.UUID `json:"follower_id" db:"follower_id"`
	FollowingID uuid.UUID `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Block struct {
	ID        uuid.UUID `json:"id" db:"id"`
	BlockerID uuid.UUID `json:"blocker_id" db:"blocker_id"`
	BlockedID uuid.UUID `json:"blocked_id" db:"blocked_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Mute struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MuterID   uuid.UUID `json:"muter_id" db:"muter_id"`
	MutedID   uuid.UUID `json:"muted_id" db:"muted_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
