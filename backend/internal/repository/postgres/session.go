package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type SessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{pool: pool}
}

func (r *SessionRepo) Create(ctx context.Context, session *model.Session) error {
	query := `INSERT INTO sessions (id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, session.ID, session.UserID, session.RefreshToken, session.UserAgent, session.IPAddress, session.ExpiresAt, session.CreatedAt)
	return err
}

func (r *SessionRepo) GetByToken(ctx context.Context, token string) (*model.Session, error) {
	s := &model.Session{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at FROM sessions WHERE refresh_token = $1`, token,
	).Scan(&s.ID, &s.UserID, &s.RefreshToken, &s.UserAgent, &s.IPAddress, &s.ExpiresAt, &s.CreatedAt)
	return s, err
}

func (r *SessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}

func (r *SessionRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}
