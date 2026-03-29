package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type MessageRepo struct {
	pool *pgxpool.Pool
}

func NewMessageRepo(pool *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{pool: pool}
}

func (r *MessageRepo) Create(ctx context.Context, msg *model.Message) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO messages (id, conversation_id, sender_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		msg.ID, msg.ConversationID, msg.SenderID, msg.Content, msg.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Update conversation timestamp
	_, err = tx.Exec(ctx,
		`UPDATE conversations SET updated_at = $2 WHERE id = $1`,
		msg.ConversationID, msg.CreatedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *MessageRepo) GetByConversation(ctx context.Context, conversationID uuid.UUID, limit int, cursor string) ([]*model.Message, error) {
	var cursorTime time.Time
	if cursor != "" {
		cursorTime, _ = time.Parse(time.RFC3339Nano, cursor)
	} else {
		cursorTime = time.Now().Add(time.Second)
	}

	query := `
		SELECT m.id, m.conversation_id, m.sender_id, m.content, m.is_deleted, m.created_at,
		       u.id, u.username, u.display_name, u.avatar_url
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = $1 AND m.created_at < $2 AND m.is_deleted = FALSE
		ORDER BY m.created_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, conversationID, cursorTime, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		m := &model.Message{Sender: &model.User{}}
		if err := rows.Scan(
			&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.IsDeleted, &m.CreatedAt,
			&m.Sender.ID, &m.Sender.Username, &m.Sender.DisplayName, &m.Sender.AvatarURL,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *MessageRepo) GetLastMessage(ctx context.Context, conversationID uuid.UUID) (*model.Message, error) {
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, m.content, m.is_deleted, m.created_at,
		       u.id, u.username, u.display_name, u.avatar_url
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = $1 AND m.is_deleted = FALSE
		ORDER BY m.created_at DESC
		LIMIT 1`

	m := &model.Message{Sender: &model.User{}}
	err := r.pool.QueryRow(ctx, query, conversationID).Scan(
		&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.IsDeleted, &m.CreatedAt,
		&m.Sender.ID, &m.Sender.Username, &m.Sender.DisplayName, &m.Sender.AvatarURL,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}
