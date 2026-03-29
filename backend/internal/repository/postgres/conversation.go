package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type ConversationRepo struct {
	pool *pgxpool.Pool
}

func NewConversationRepo(pool *pgxpool.Pool) *ConversationRepo {
	return &ConversationRepo{pool: pool}
}

func (r *ConversationRepo) Create(ctx context.Context, conv *model.Conversation, participantIDs []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO conversations (id, is_group, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		conv.ID, conv.IsGroup, conv.Name, conv.CreatedAt, conv.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, uid := range participantIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO conversation_participants (id, conversation_id, user_id, last_read_at, joined_at) VALUES ($1, $2, $3, $4, $5)`,
			uuid.New(), conv.ID, uid, conv.CreatedAt, conv.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *ConversationRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Conversation, error) {
	conv := &model.Conversation{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, is_group, name, created_at, updated_at FROM conversations WHERE id = $1`, id,
	).Scan(&conv.ID, &conv.IsGroup, &conv.Name, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	participants, err := r.getParticipants(ctx, id)
	if err != nil {
		return nil, err
	}
	conv.Participants = participants

	return conv, nil
}

func (r *ConversationRepo) GetByUserID(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Conversation, error) {
	var cursorTime time.Time
	if cursor != "" {
		cursorTime, _ = time.Parse(time.RFC3339Nano, cursor)
	} else {
		cursorTime = time.Now().Add(time.Second)
	}

	query := `
		SELECT c.id, c.is_group, c.name, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_participants cp ON cp.conversation_id = c.id
		WHERE cp.user_id = $1 AND c.updated_at < $2
		ORDER BY c.updated_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, userID, cursorTime, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []*model.Conversation
	for rows.Next() {
		c := &model.Conversation{}
		if err := rows.Scan(&c.ID, &c.IsGroup, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}

	// Load participants for each conversation
	for _, c := range convs {
		participants, err := r.getParticipants(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		c.Participants = participants
	}

	return convs, nil
}

func (r *ConversationRepo) GetDMBetween(ctx context.Context, userID1, userID2 uuid.UUID) (*model.Conversation, error) {
	query := `
		SELECT c.id, c.is_group, c.name, c.created_at, c.updated_at
		FROM conversations c
		WHERE c.is_group = FALSE
		  AND c.id IN (SELECT conversation_id FROM conversation_participants WHERE user_id = $1)
		  AND c.id IN (SELECT conversation_id FROM conversation_participants WHERE user_id = $2)
		LIMIT 1`

	conv := &model.Conversation{}
	err := r.pool.QueryRow(ctx, query, userID1, userID2).Scan(&conv.ID, &conv.IsGroup, &conv.Name, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return conv, nil
}

func (r *ConversationRepo) MarkRead(ctx context.Context, conversationID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE conversation_participants SET last_read_at = NOW() WHERE conversation_id = $1 AND user_id = $2`,
		conversationID, userID,
	)
	return err
}

func (r *ConversationRepo) UnreadCount(ctx context.Context, conversationID, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM messages m
		JOIN conversation_participants cp ON cp.conversation_id = m.conversation_id AND cp.user_id = $2
		WHERE m.conversation_id = $1 AND m.created_at > cp.last_read_at AND m.sender_id != $2`,
		conversationID, userID,
	).Scan(&count)
	return count, err
}

func (r *ConversationRepo) getParticipants(ctx context.Context, convID uuid.UUID) ([]model.User, error) {
	query := `
		SELECT u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM conversation_participants cp
		JOIN users u ON u.id = cp.user_id
		WHERE cp.conversation_id = $1`

	rows, err := r.pool.Query(ctx, query, convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		u := model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.IsVerified); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
