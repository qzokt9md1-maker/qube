package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/model"
	"github.com/kuzuokatakumi/qube/internal/repository/postgres"
	"github.com/kuzuokatakumi/qube/internal/ws"
)

type NotificationService struct {
	notifRepo *postgres.NotificationRepo
	hub       *ws.Hub
}

func NewNotificationService(notifRepo *postgres.NotificationRepo, hub *ws.Hub) *NotificationService {
	return &NotificationService{
		notifRepo: notifRepo,
		hub:       hub,
	}
}

func (s *NotificationService) Create(ctx context.Context, userID, actorID uuid.UUID, notifType string, postID *uuid.UUID) {
	if userID == actorID {
		return
	}

	notif := &model.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		ActorID:   actorID,
		Type:      notifType,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	if err := s.notifRepo.Create(ctx, notif); err != nil {
		return
	}

	// Real-time push via WebSocket
	if s.hub != nil {
		s.hub.SendToUser(userID, ws.Event{
			Type: "notification",
			Payload: map[string]interface{}{
				"id":   notif.ID.String(),
				"type": notifType,
			},
		})
	}
}

func (s *NotificationService) GetByUserID(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Notification, int, error) {
	notifs, err := s.notifRepo.GetByUserID(ctx, userID, limit, cursor)
	if err != nil {
		return nil, 0, err
	}
	unread, _ := s.notifRepo.UnreadCount(ctx, userID)
	return notifs, unread, nil
}

func (s *NotificationService) MarkRead(ctx context.Context, ids []uuid.UUID) error {
	return s.notifRepo.MarkRead(ctx, ids)
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return s.notifRepo.MarkAllRead(ctx, userID)
}
