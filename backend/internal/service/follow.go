package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/model"
	"github.com/kuzuokatakumi/qube/internal/repository/postgres"
)

type FollowService struct {
	followRepo   *postgres.FollowRepo
	userRepo     *postgres.UserRepo
	blockRepo    *postgres.BlockRepo
	notifService *NotificationService
}

func NewFollowService(
	followRepo *postgres.FollowRepo,
	userRepo *postgres.UserRepo,
	blockRepo *postgres.BlockRepo,
	notifService *NotificationService,
) *FollowService {
	return &FollowService{
		followRepo:   followRepo,
		userRepo:     userRepo,
		blockRepo:    blockRepo,
		notifService: notifService,
	}
}

func (s *FollowService) Follow(ctx context.Context, followerID, followingID uuid.UUID) (*model.User, error) {
	if followerID == followingID {
		return nil, ErrSelfAction
	}

	blocked, _ := s.blockRepo.IsBlocked(ctx, followingID, followerID)
	if blocked {
		return nil, ErrBlocked
	}

	follow := &model.Follow{
		ID:          uuid.New(),
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}

	if err := s.followRepo.Create(ctx, follow); err != nil {
		return nil, err
	}

	_ = s.userRepo.UpdateCounts(ctx, followerID, "following_count", 1)
	_ = s.userRepo.UpdateCounts(ctx, followingID, "follower_count", 1)

	// Notify
	s.notifService.Create(ctx, followingID, followerID, "follow", nil)

	return s.userRepo.GetByID(ctx, followingID)
}

func (s *FollowService) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) (*model.User, error) {
	if err := s.followRepo.Delete(ctx, followerID, followingID); err != nil {
		return nil, err
	}

	_ = s.userRepo.UpdateCounts(ctx, followerID, "following_count", -1)
	_ = s.userRepo.UpdateCounts(ctx, followingID, "follower_count", -1)

	return s.userRepo.GetByID(ctx, followingID)
}

func (s *FollowService) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return s.followRepo.Exists(ctx, followerID, followingID)
}

func (s *FollowService) GetFollowers(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.User, error) {
	return s.followRepo.GetFollowers(ctx, userID, limit, cursor)
}

func (s *FollowService) GetFollowing(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.User, error) {
	return s.followRepo.GetFollowing(ctx, userID, limit, cursor)
}
