package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type PostRepo struct {
	pool *pgxpool.Pool
}

func NewPostRepo(pool *pgxpool.Pool) *PostRepo {
	return &PostRepo{pool: pool}
}

func (r *PostRepo) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO posts (id, user_id, content, reply_to_id, repost_of_id, quote_of_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query,
		post.ID, post.UserID, post.Content, post.ReplyToID, post.RepostOfID, post.QuoteOfID,
		post.CreatedAt, post.UpdatedAt,
	)
	return err
}

func (r *PostRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.content, p.reply_to_id, p.repost_of_id, p.quote_of_id,
		       p.like_count, p.repost_count, p.reply_count, p.quote_count, p.is_deleted, p.created_at, p.updated_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = $1 AND p.is_deleted = FALSE`
	p := &model.Post{User: &model.User{}}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.UserID, &p.Content, &p.ReplyToID, &p.RepostOfID, &p.QuoteOfID,
		&p.LikeCount, &p.RepostCount, &p.ReplyCount, &p.QuoteCount, &p.IsDeleted, &p.CreatedAt, &p.UpdatedAt,
		&p.User.ID, &p.User.Username, &p.User.DisplayName, &p.User.AvatarURL, &p.User.IsVerified,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE posts SET is_deleted = TRUE WHERE id = $1`, id)
	return err
}

func (r *PostRepo) GetTimeline(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error) {
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
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.user_id IN (SELECT following_id FROM follows WHERE follower_id = $1)
		  AND p.is_deleted = FALSE
		  AND p.created_at < $2
		ORDER BY p.created_at DESC
		LIMIT $3`

	return r.scanPosts(ctx, query, userID, cursorTime, limit)
}

func (r *PostRepo) GetUserPosts(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error) {
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
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.user_id = $1 AND p.is_deleted = FALSE AND p.created_at < $2 AND p.reply_to_id IS NULL
		ORDER BY p.created_at DESC
		LIMIT $3`

	return r.scanPosts(ctx, query, userID, cursorTime, limit)
}

func (r *PostRepo) GetReplies(ctx context.Context, postID uuid.UUID, limit int, cursor string) ([]*model.Post, error) {
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
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.reply_to_id = $1 AND p.is_deleted = FALSE AND p.created_at < $2
		ORDER BY p.created_at ASC
		LIMIT $3`

	return r.scanPosts(ctx, query, postID, cursorTime, limit)
}

func (r *PostRepo) GetMediaByPostIDs(ctx context.Context, postIDs []uuid.UUID) (map[uuid.UUID][]model.Media, error) {
	query := `SELECT id, post_id, media_type, url, thumbnail_url, width, height, sort_order, created_at FROM post_media WHERE post_id = ANY($1) ORDER BY sort_order`
	rows, err := r.pool.Query(ctx, query, postIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]model.Media)
	for rows.Next() {
		m := model.Media{}
		if err := rows.Scan(&m.ID, &m.PostID, &m.MediaType, &m.URL, &m.ThumbnailURL, &m.Width, &m.Height, &m.SortOrder, &m.CreatedAt); err != nil {
			return nil, err
		}
		result[m.PostID] = append(result[m.PostID], m)
	}
	return result, nil
}

func (r *PostRepo) CreateMedia(ctx context.Context, media *model.Media) error {
	query := `INSERT INTO post_media (id, post_id, media_type, url, thumbnail_url, width, height, sort_order, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query, media.ID, media.PostID, media.MediaType, media.URL, media.ThumbnailURL, media.Width, media.Height, media.SortOrder, media.CreatedAt)
	return err
}

func (r *PostRepo) UpdateCounts(ctx context.Context, id uuid.UUID, field string, delta int) error {
	allowed := map[string]bool{"like_count": true, "repost_count": true, "reply_count": true, "quote_count": true}
	if !allowed[field] {
		return nil
	}
	query := `UPDATE posts SET ` + field + ` = GREATEST(` + field + ` + $2, 0) WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, delta)
	return err
}

func (r *PostRepo) scanPosts(ctx context.Context, query string, args ...interface{}) ([]*model.Post, error) {
	rows, err := r.pool.Query(ctx, query, args...)
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
