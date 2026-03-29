package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type FollowRepo struct {
	pool *pgxpool.Pool
}

func NewFollowRepo(pool *pgxpool.Pool) *FollowRepo {
	return &FollowRepo{pool: pool}
}

func (r *FollowRepo) Create(ctx context.Context, follow *model.Follow) error {
	query := `INSERT INTO follows (id, follower_id, following_id, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, follow.ID, follow.FollowerID, follow.FollowingID, follow.CreatedAt)
	return err
}

func (r *FollowRepo) Delete(ctx context.Context, followerID, followingID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`, followerID, followingID)
	return err
}

func (r *FollowRepo) Exists(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)`, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *FollowRepo) GetFollowers(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.User, error) {
	var cursorTime time.Time
	if cursor != "" {
		cursorTime, _ = time.Parse(time.RFC3339Nano, cursor)
	} else {
		cursorTime = time.Now().Add(time.Second)
	}

	query := `
		SELECT u.id, u.username, u.display_name, u.email, u.password_hash, u.bio, u.avatar_url, u.header_url, u.location, u.website, u.is_verified, u.is_private, u.follower_count, u.following_count, u.post_count, u.created_at, u.updated_at
		FROM follows f
		JOIN users u ON u.id = f.follower_id
		WHERE f.following_id = $1 AND f.created_at < $2
		ORDER BY f.created_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, userID, cursorTime, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}

func (r *FollowRepo) GetFollowing(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.User, error) {
	var cursorTime time.Time
	if cursor != "" {
		cursorTime, _ = time.Parse(time.RFC3339Nano, cursor)
	} else {
		cursorTime = time.Now().Add(time.Second)
	}

	query := `
		SELECT u.id, u.username, u.display_name, u.email, u.password_hash, u.bio, u.avatar_url, u.header_url, u.location, u.website, u.is_verified, u.is_private, u.follower_count, u.following_count, u.post_count, u.created_at, u.updated_at
		FROM follows f
		JOIN users u ON u.id = f.following_id
		WHERE f.follower_id = $1 AND f.created_at < $2
		ORDER BY f.created_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, userID, cursorTime, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}

func (r *FollowRepo) GetFollowingIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, `SELECT following_id FROM follows WHERE follower_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
