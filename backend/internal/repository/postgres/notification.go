package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kuzuokatakumi/qube/internal/model"
)

type NotificationRepo struct {
	pool *pgxpool.Pool
}

func NewNotificationRepo(pool *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{pool: pool}
}

func (r *NotificationRepo) Create(ctx context.Context, notif *model.Notification) error {
	query := `INSERT INTO notifications (id, user_id, actor_id, type, post_id, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, notif.ID, notif.UserID, notif.ActorID, notif.Type, notif.PostID, notif.CreatedAt)
	return err
}

func (r *NotificationRepo) GetByUserID(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Notification, error) {
	var cursorTime time.Time
	if cursor != "" {
		cursorTime, _ = time.Parse(time.RFC3339Nano, cursor)
	} else {
		cursorTime = time.Now().Add(time.Second)
	}

	query := `
		SELECT n.id, n.user_id, n.actor_id, n.type, n.post_id, n.is_read, n.created_at,
		       u.id, u.username, u.display_name, u.avatar_url, u.is_verified
		FROM notifications n
		JOIN users u ON u.id = n.actor_id
		WHERE n.user_id = $1 AND n.created_at < $2
		ORDER BY n.created_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, userID, cursorTime, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []*model.Notification
	for rows.Next() {
		n := &model.Notification{Actor: &model.User{}}
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.ActorID, &n.Type, &n.PostID, &n.IsRead, &n.CreatedAt,
			&n.Actor.ID, &n.Actor.Username, &n.Actor.DisplayName, &n.Actor.AvatarURL, &n.Actor.IsVerified,
		); err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (r *NotificationRepo) MarkRead(ctx context.Context, ids []uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET is_read = TRUE WHERE id = ANY($1)`, ids)
	return err
}

func (r *NotificationRepo) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`, userID)
	return err
}

func (r *NotificationRepo) UnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`, userID).Scan(&count)
	return count, err
}
