package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlockRepo struct {
	pool *pgxpool.Pool
}

func NewBlockRepo(pool *pgxpool.Pool) *BlockRepo {
	return &BlockRepo{pool: pool}
}

func (r *BlockRepo) Block(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO blocks (id, blocker_id, blocked_id, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		uuid.New(), blockerID, blockedID, time.Now(),
	)
	return err
}

func (r *BlockRepo) Unblock(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM blocks WHERE blocker_id = $1 AND blocked_id = $2`, blockerID, blockedID)
	return err
}

func (r *BlockRepo) IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM blocks WHERE blocker_id = $1 AND blocked_id = $2)`, blockerID, blockedID).Scan(&exists)
	return exists, err
}

type MuteRepo struct {
	pool *pgxpool.Pool
}

func NewMuteRepo(pool *pgxpool.Pool) *MuteRepo {
	return &MuteRepo{pool: pool}
}

func (r *MuteRepo) Mute(ctx context.Context, muterID, mutedID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO mutes (id, muter_id, muted_id, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		uuid.New(), muterID, mutedID, time.Now(),
	)
	return err
}

func (r *MuteRepo) Unmute(ctx context.Context, muterID, mutedID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM mutes WHERE muter_id = $1 AND muted_id = $2`, muterID, mutedID)
	return err
}

func (r *MuteRepo) IsMuted(ctx context.Context, muterID, mutedID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM mutes WHERE muter_id = $1 AND muted_id = $2)`, muterID, mutedID).Scan(&exists)
	return exists, err
}
