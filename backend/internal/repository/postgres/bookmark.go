package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type BookmarkRepo struct {
	pool *pgxpool.Pool
}

func NewBookmarkRepo(pool *pgxpool.Pool) *BookmarkRepo {
	return &BookmarkRepo{pool: pool}
}

func (r *BookmarkRepo) Create(ctx context.Context, userID, postID uuid.UUID) error {
	query := `INSERT INTO bookmarks (id, user_id, post_id, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, query, uuid.New(), userID, postID, time.Now())
	return err
}

func (r *BookmarkRepo) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM bookmarks WHERE user_id = $1 AND post_id = $2`, userID, postID)
	return err
}

func (r *BookmarkRepo) Exists(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND post_id = $2)`, userID, postID).Scan(&exists)
	return exists, err
}

func (r *BookmarkRepo) GetByUserID(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error) {
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
		FROM bookmarks b
		JOIN posts p ON p.id = b.post_id
		JOIN users u ON u.id = p.user_id
		WHERE b.user_id = $1 AND b.created_at < $2 AND p.is_deleted = FALSE
		ORDER BY b.created_at DESC
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
