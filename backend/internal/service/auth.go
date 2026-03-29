package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/config"
	"github.com/kuzuokatakumi/qube/internal/model"
	"github.com/kuzuokatakumi/qube/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already taken")
	ErrUsernameTaken      = errors.New("username already taken")
)

type AuthService struct {
	userRepo   repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtConfig  config.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, jwtCfg config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		sessionRepo: sessionRepo,
		jwtConfig:  jwtCfg,
	}
}

type AuthPayload struct {
	AccessToken  string
	RefreshToken string
	User         *model.User
}

func (s *AuthService) Register(ctx context.Context, username, displayName, email, password string) (*AuthPayload, error) {
	// Check uniqueness
	if exists, _ := s.userRepo.ExistsByEmail(ctx, email); exists {
		return nil, ErrEmailTaken
	}
	if exists, _ := s.userRepo.ExistsByUsername(ctx, username); exists {
		return nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:          uuid.New(),
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		PasswordHash: string(hash),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthPayload, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthPayload, error) {
	session, err := s.sessionRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.sessionRepo.Delete(ctx, session.ID)
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	// Rotate refresh token
	_ = s.sessionRepo.Delete(ctx, session.ID)
	return s.generateTokens(ctx, user)
}

func (s *AuthService) generateTokens(ctx context.Context, user *model.User) (*AuthPayload, error) {
	now := time.Now()

	// Access token
	accessClaims := jwt.MapClaims{
		"sub": user.ID.String(),
		"iat": now.Unix(),
		"exp": now.Add(s.jwtConfig.AccessTokenTTL).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshBytes := make([]byte, 32)
	if _, err := rand.Read(refreshBytes); err != nil {
		return nil, err
	}
	refreshStr := hex.EncodeToString(refreshBytes)

	session := &model.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshStr,
		ExpiresAt:    now.Add(s.jwtConfig.RefreshTokenTTL),
		CreatedAt:    now,
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &AuthPayload{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		User:         user,
	}, nil
}
