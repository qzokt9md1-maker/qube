package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/model"
	"github.com/kuzuokatakumi/qube/internal/repository/postgres"
	"github.com/kuzuokatakumi/qube/internal/ws"
)

type DMService struct {
	convRepo    *postgres.ConversationRepo
	msgRepo     *postgres.MessageRepo
	blockRepo   *postgres.BlockRepo
	notifService *NotificationService
	hub         *ws.Hub
}

func NewDMService(
	convRepo *postgres.ConversationRepo,
	msgRepo *postgres.MessageRepo,
	blockRepo *postgres.BlockRepo,
	notifService *NotificationService,
	hub *ws.Hub,
) *DMService {
	return &DMService{
		convRepo:    convRepo,
		msgRepo:     msgRepo,
		blockRepo:   blockRepo,
		notifService: notifService,
		hub:         hub,
	}
}

func (s *DMService) CreateConversation(ctx context.Context, creatorID uuid.UUID, participantIDs []uuid.UUID, firstMessage string) (*model.Conversation, error) {
	allIDs := append([]uuid.UUID{creatorID}, participantIDs...)

	// For 1-on-1, check if conversation already exists
	if len(participantIDs) == 1 {
		existing, err := s.convRepo.GetDMBetween(ctx, creatorID, participantIDs[0])
		if err == nil && existing != nil {
			// Send message to existing conversation
			_, err := s.SendMessage(ctx, creatorID, existing.ID, firstMessage)
			if err != nil {
				return nil, err
			}
			return s.convRepo.GetByID(ctx, existing.ID)
		}
	}

	now := time.Now()
	conv := &model.Conversation{
		ID:        uuid.New(),
		IsGroup:   len(participantIDs) > 1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.convRepo.Create(ctx, conv, allIDs); err != nil {
		return nil, err
	}

	// Send first message
	if firstMessage != "" {
		_, err := s.SendMessage(ctx, creatorID, conv.ID, firstMessage)
		if err != nil {
			return nil, err
		}
	}

	return s.convRepo.GetByID(ctx, conv.ID)
}

func (s *DMService) SendMessage(ctx context.Context, senderID, conversationID uuid.UUID, content string) (*model.Message, error) {
	now := time.Now()
	msg := &model.Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      now,
	}

	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	// Real-time push via WebSocket
	if s.hub != nil {
		conv, err := s.convRepo.GetByID(ctx, conversationID)
		if err == nil {
			for _, p := range conv.Participants {
				if p.ID != senderID {
					s.hub.SendToUser(p.ID, ws.Event{
						Type: "new_message",
						Payload: map[string]interface{}{
							"conversation_id": conversationID.String(),
							"message_id":      msg.ID.String(),
							"sender_id":       senderID.String(),
							"content":         content,
						},
					})
					// DM notification
					s.notifService.Create(ctx, p.ID, senderID, "dm", nil)
				}
			}
		}
	}

	return msg, nil
}

func (s *DMService) GetConversations(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Conversation, error) {
	convs, err := s.convRepo.GetByUserID(ctx, userID, limit, cursor)
	if err != nil {
		return nil, err
	}

	// Load last message for each
	for _, c := range convs {
		lastMsg, err := s.msgRepo.GetLastMessage(ctx, c.ID)
		if err == nil {
			c.LastMessage = lastMsg
		}
	}

	return convs, nil
}

func (s *DMService) GetMessages(ctx context.Context, conversationID uuid.UUID, limit int, cursor string) ([]*model.Message, error) {
	return s.msgRepo.GetByConversation(ctx, conversationID, limit, cursor)
}

func (s *DMService) MarkRead(ctx context.Context, conversationID, userID uuid.UUID) error {
	return s.convRepo.MarkRead(ctx, conversationID, userID)
}
