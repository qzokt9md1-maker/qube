package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type LikeRepo struct {
	pool *pgxpool.Pool
}

func NewLikeRepo(pool *pgxpool.Pool) *LikeRepo {
	return &LikeRepo{pool: pool}
}

func (r *LikeRepo) Create(ctx context.Context, userID, postID uuid.UUID) error {
	query := `INSERT INTO likes (id, user_id, post_id, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, uuid.New(), userID, postID, time.Now())
	return err
}

func (r *LikeRepo) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM likes WHERE user_id = $1 AND post_id = $2`, userID, postID)
	return err
}

func (r *LikeRepo) Exists(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)`, userID, postID).Scan(&exists)
	return exists, err
}

func (r *LikeRepo) GetUserLikes(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error) {
	var cursorTime time.Time
	if cursor != "" {
		cursorTime, _ = time.Parse(time.RFC3339Nano, cursor)
	} else {
		cursorTime = time.Now().Add(time.Second)
	}

	query := `
		SELECT p.id, p.user_id, p.content, p.reply_to_id, p.repost_of_id, p.quote_of_id,
		       p.like_count, p.repost_count, p.reply_count, p.quote_count, p.is_deleted, p.created_at, p.updated_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM likes l
		JOIN posts p ON p.id = l.post_id
		JOIN users u ON u.id = p.user_id
		WHERE l.user_id = $1 AND l.created_at < $2 AND p.is_deleted = FALSE
		ORDER BY l.created_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, userID, cursorTime, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		p := &model.Post{User: &model.User{}}
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Content, &p.ReplyToID, &p.RepostOfID, &p.QuoteOfID,
			&p.LikeCount, &p.RepostCount, &p.ReplyCount, &p.QuoteCount, &p.IsDeleted, &p.CreatedAt, &p.UpdatedAt,
			&p.User.ID, &p.User.Username, &p.User.DisplayName, &p.User.AvatarURL, &p.User.IsVerified,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *LikeRepo) IsLikedByUser(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	query := `SELECT post_id FROM likes WHERE user_id = $1 AND post_id = ANY($2)`
	rows, err := r.pool.Query(ctx, query, userID, postIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]bool)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result[id] = true
	}
	return result, nil
}
