package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/model"
	"github.com/kuzuokatakumi/qube/internal/repository/postgres"
)

type UserService struct {
	userRepo  *postgres.UserRepo
	blockRepo *postgres.BlockRepo
	muteRepo  *postgres.MuteRepo
}

func NewUserService(userRepo *postgres.UserRepo, blockRepo *postgres.BlockRepo, muteRepo *postgres.MuteRepo) *UserService {
	return &UserService{
		userRepo:  userRepo,
		blockRepo: blockRepo,
		muteRepo:  muteRepo,
	}
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, displayName, bio, location, website *string, isPrivate *bool) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if displayName != nil {
		user.DisplayName = *displayName
	}
	if bio != nil {
		user.Bio = *bio
	}
	if location != nil {
		user.Location = *location
	}
	if website != nil {
		user.Website = *website
	}
	if isPrivate != nil {
		user.IsPrivate = *isPrivate
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.AvatarURL = avatarURL
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Search(ctx context.Context, query string, limit int, cursor string) ([]*model.User, error) {
	return s.userRepo.Search(ctx, query, limit, cursor)
}

func (s *UserService) Block(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	if blockerID == blockedID {
		return ErrSelfAction
	}
	return s.blockRepo.Block(ctx, blockerID, blockedID)
}

func (s *UserService) Unblock(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	return s.blockRepo.Unblock(ctx, blockerID, blockedID)
}

func (s *UserService) Mute(ctx context.Context, muterID, mutedID uuid.UUID) error {
	if muterID == mutedID {
		return ErrSelfAction
	}
	return s.muteRepo.Mute(ctx, muterID, mutedID)
}

func (s *UserService) Unmute(ctx context.Context, muterID, mutedID uuid.UUID) error {
	return s.muteRepo.Unmute(ctx, muterID, mutedID)
}

func (s *UserService) IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
	return s.blockRepo.IsBlocked(ctx, blockerID, blockedID)
}

func (s *UserService) IsMuted(ctx context.Context, muterID, mutedID uuid.UUID) (bool, error) {
	return s.muteRepo.IsMuted(ctx, muterID, mutedID)
}
