package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, username, display_name, email, password_hash, bio, avatar_url, header_url, location, website, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Username, user.DisplayName, user.Email, user.PasswordHash,
		user.Bio, user.AvatarURL, user.HeaderURL, user.Location, user.Website,
		user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `SELECT id, username, display_name, email, password_hash, bio, avatar_url, header_url, location, website, is_verified, is_private, follower_count, following_count, post_count, created_at, updated_at FROM users WHERE id = $1`
	return r.scanUser(ctx, query, id)
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, username, display_name, email, password_hash, bio, avatar_url, header_url, location, website, is_verified, is_private, follower_count, following_count, post_count, created_at, updated_at FROM users WHERE username = $1`
	return r.scanUser(ctx, query, username)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, username, display_name, email, password_hash, bio, avatar_url, header_url, location, website, is_verified, is_private, follower_count, following_count, post_count, created_at, updated_at FROM users WHERE email = $1`
	return r.scanUser(ctx, query, email)
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}

func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username).Scan(&exists)
	return exists, err
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users SET display_name = $2, bio = $3, avatar_url = $4, header_url = $5, location = $6, website = $7, is_private = $8
		WHERE id = $1`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.DisplayName, user.Bio, user.AvatarURL, user.HeaderURL,
		user.Location, user.Website, user.IsPrivate,
	)
	return err
}

func (r *UserRepo) UpdateCounts(ctx context.Context, id uuid.UUID, field string, delta int) error {
	allowed := map[string]bool{"follower_count": true, "following_count": true, "post_count": true}
	if !allowed[field] {
		return fmt.Errorf("invalid count field: %s", field)
	}
	query := fmt.Sprintf(`UPDATE users SET %s = GREATEST(%s + $2, 0) WHERE id = $1`, field, field)
	_, err := r.pool.Exec(ctx, query, id, delta)
	return err
}

func (r *UserRepo) Search(ctx context.Context, query string, limit int, cursor string) ([]*model.User, error) {
	pattern := "%" + query + "%"
	sql := `
		SELECT id, username, display_name, email, password_hash, bio, avatar_url, header_url, location, website, is_verified, is_private, follower_count, following_count, post_count, created_at, updated_at
		FROM users
		WHERE (username ILIKE $1 OR display_name ILIKE $1)
		ORDER BY follower_count DESC
		LIMIT $2`
	rows, err := r.pool.Query(ctx, sql, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.HeaderURL, &u.Location, &u.Website, &u.IsVerified, &u.IsPrivate, &u.FollowerCount, &u.FollowingCount, &u.PostCount, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.User, error) {
	query := `
		SELECT id, username, display_name, email, password_hash, bio, avatar_url, header_url, location, website, is_verified, is_private, follower_count, following_count, post_count, created_at, updated_at
		FROM users WHERE id = ANY($1)`
	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.HeaderURL, &u.Location, &u.Website, &u.IsVerified, &u.IsPrivate, &u.FollowerCount, &u.FollowingCount, &u.PostCount, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) scanUser(ctx context.Context, query string, arg interface{}) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx, query, arg).Scan(
		&u.ID, &u.Username, &u.DisplayName, &u.Email, &u.PasswordHash,
		&u.Bio, &u.AvatarURL, &u.HeaderURL, &u.Location, &u.Website,
		&u.IsVerified, &u.IsPrivate, &u.FollowerCount, &u.FollowingCount, &u.PostCount,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}
