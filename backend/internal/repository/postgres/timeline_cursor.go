package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimelineCursorRepo struct {
	pool *pgxpool.Pool
}

func NewTimelineCursorRepo(pool *pgxpool.Pool) *TimelineCursorRepo {
	return &TimelineCursorRepo{pool: pool}
}

func (r *TimelineCursorRepo) Get(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	var postID *uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT last_seen_post_id FROM timeline_cursors WHERE user_id = $1`, userID,
	).Scan(&postID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return postID, nil
}

func (r *TimelineCursorRepo) Update(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO timeline_cursors (user_id, last_seen_at, last_seen_post_id)
		VALUES ($1, NOW(), $2)
		ON CONFLICT (user_id)
		DO UPDATE SET last_seen_at = NOW(), last_seen_post_id = $2`,
		userID, postID,
	)
	return err
}

func (r *TimelineCursorRepo) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM posts p
		JOIN follows f ON f.following_id = p.user_id AND f.follower_id = $1
		LEFT JOIN timeline_cursors tc ON tc.user_id = $1
		WHERE p.is_deleted = FALSE
		  AND (tc.last_seen_at IS NULL OR p.created_at > tc.last_seen_at)`,
		userID,
	).Scan(&count)
	return count, err
}
