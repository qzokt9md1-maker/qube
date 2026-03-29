package model

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID `json:"id" db:"id"`
	IsGroup   bool      `json:"is_group" db:"is_group"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	// Joined fields
	Participants []User   `json:"participants,omitempty"`
	LastMessage  *Message `json:"last_message,omitempty"`
}

type Message struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id" db:"sender_id"`
	Content        string    `json:"content" db:"content"`
	IsDeleted      bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	// Joined fields
	Sender *User `json:"sender,omitempty"`
}

type ConversationParticipant struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	LastReadAt     time.Time `json:"last_read_at" db:"last_read_at"`
	JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
}
