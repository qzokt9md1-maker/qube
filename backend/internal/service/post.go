package service

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/model"
	"github.com/kuzuokatakumi/qube/internal/repository/postgres"
)

var hashtagRegex = regexp.MustCompile(`#(\w+)`)

type PostService struct {
	postRepo     *postgres.PostRepo
	userRepo     *postgres.UserRepo
	hashtagRepo  *postgres.HashtagRepo
	likeRepo     *postgres.LikeRepo
	bookmarkRepo *postgres.BookmarkRepo
	notifService *NotificationService
	timelineSvc  *TimelineService
}

func NewPostService(
	postRepo *postgres.PostRepo,
	userRepo *postgres.UserRepo,
	hashtagRepo *postgres.HashtagRepo,
	likeRepo *postgres.LikeRepo,
	bookmarkRepo *postgres.BookmarkRepo,
	notifService *NotificationService,
	timelineSvc *TimelineService,
) *PostService {
	return &PostService{
		postRepo:     postRepo,
		userRepo:     userRepo,
		hashtagRepo:  hashtagRepo,
		likeRepo:     likeRepo,
		bookmarkRepo: bookmarkRepo,
		notifService: notifService,
		timelineSvc:  timelineSvc,
	}
}

func (s *PostService) Create(ctx context.Context, userID uuid.UUID, content string, replyToID, quoteOfID *uuid.UUID, mediaIDs []uuid.UUID) (*model.Post, error) {
	now := time.Now()
	post := &model.Post{
		ID:        uuid.New(),
		UserID:    userID,
		Content:   content,
		ReplyToID: replyToID,
		QuoteOfID: quoteOfID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	// Update user post count
	_ = s.userRepo.UpdateCounts(ctx, userID, "post_count", 1)

	// Extract and link hashtags
	tags := extractHashtags(content)
	if len(tags) > 0 {
		_ = s.hashtagRepo.UpsertAndLink(ctx, post.ID, tags)
	}

	// Update reply count on parent
	if replyToID != nil {
		_ = s.postRepo.UpdateCounts(ctx, *replyToID, "reply_count", 1)
		// Notify parent post author
		parent, err := s.postRepo.GetByID(ctx, *replyToID)
		if err == nil && parent.UserID != userID {
			s.notifService.Create(ctx, parent.UserID, userID, "reply", &post.ID)
		}
	}

	// Update quote count
	if quoteOfID != nil {
		_ = s.postRepo.UpdateCounts(ctx, *quoteOfID, "quote_count", 1)
		quoted, err := s.postRepo.GetByID(ctx, *quoteOfID)
		if err == nil && quoted.UserID != userID {
			s.notifService.Create(ctx, quoted.UserID, userID, "quote", &post.ID)
		}
	}

	// Extract mentions and notify
	mentions := extractMentions(content)
	for _, username := range mentions {
		mentioned, err := s.userRepo.GetByUsername(ctx, username)
		if err == nil && mentioned.ID != userID {
			s.notifService.Create(ctx, mentioned.ID, userID, "mention", &post.ID)
		}
	}

	// Fan-out to timeline via Redis
	if s.timelineSvc != nil {
		go s.timelineSvc.FanOutPost(context.Background(), post)
	}

	// Reload with user
	return s.postRepo.GetByID(ctx, post.ID)
}

func (s *PostService) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return ErrUnauthorizedAction
	}

	if err := s.postRepo.Delete(ctx, postID); err != nil {
		return err
	}

	_ = s.userRepo.UpdateCounts(ctx, userID, "post_count", -1)

	if post.ReplyToID != nil {
		_ = s.postRepo.UpdateCounts(ctx, *post.ReplyToID, "reply_count", -1)
	}
	return nil
}

func (s *PostService) Like(ctx context.Context, userID, postID uuid.UUID) (*model.Post, error) {
	if err := s.likeRepo.Create(ctx, userID, postID); err != nil {
		return nil, err
	}
	_ = s.postRepo.UpdateCounts(ctx, postID, "like_count", 1)

	// Notify
	post, err := s.postRepo.GetByID(ctx, postID)
	if err == nil && post.UserID != userID {
		s.notifService.Create(ctx, post.UserID, userID, "like", &postID)
	}
	return post, err
}

func (s *PostService) Unlike(ctx context.Context, userID, postID uuid.UUID) (*model.Post, error) {
	if err := s.likeRepo.Delete(ctx, userID, postID); err != nil {
		return nil, err
	}
	_ = s.postRepo.UpdateCounts(ctx, postID, "like_count", -1)
	return s.postRepo.GetByID(ctx, postID)
}

func (s *PostService) Repost(ctx context.Context, userID, postID uuid.UUID) (*model.Post, error) {
	now := time.Now()
	repost := &model.Post{
		ID:         uuid.New(),
		UserID:     userID,
		Content:    "",
		RepostOfID: &postID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.postRepo.Create(ctx, repost); err != nil {
		return nil, err
	}
	_ = s.postRepo.UpdateCounts(ctx, postID, "repost_count", 1)

	original, err := s.postRepo.GetByID(ctx, postID)
	if err == nil && original.UserID != userID {
		s.notifService.Create(ctx, original.UserID, userID, "repost", &postID)
	}

	if s.timelineSvc != nil {
		go s.timelineSvc.FanOutPost(context.Background(), repost)
	}

	return s.postRepo.GetByID(ctx, repost.ID)
}

func (s *PostService) Unrepost(ctx context.Context, userID, postID uuid.UUID) error {
	_ = s.postRepo.UpdateCounts(ctx, postID, "repost_count", -1)
	return nil
}

func (s *PostService) Bookmark(ctx context.Context, userID, postID uuid.UUID) error {
	return s.bookmarkRepo.Create(ctx, userID, postID)
}

func (s *PostService) Unbookmark(ctx context.Context, userID, postID uuid.UUID) error {
	return s.bookmarkRepo.Delete(ctx, userID, postID)
}

func (s *PostService) GetByID(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	return s.postRepo.GetByID(ctx, id)
}

func (s *PostService) GetUserPosts(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error) {
	return s.postRepo.GetUserPosts(ctx, userID, limit, cursor)
}

func (s *PostService) GetReplies(ctx context.Context, postID uuid.UUID, limit int, cursor string) ([]*model.Post, error) {
	return s.postRepo.GetReplies(ctx, postID, limit, cursor)
}

func (s *PostService) GetUserLikes(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error) {
	return s.likeRepo.GetUserLikes(ctx, userID, limit, cursor)
}

func (s *PostService) GetBookmarks(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, error) {
	return s.bookmarkRepo.GetByUserID(ctx, userID, limit, cursor)
}

func extractHashtags(content string) []string {
	matches := hashtagRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var tags []string
	for _, m := range matches {
		tag := strings.ToLower(m[1])
		if !seen[tag] {
			seen[tag] = true
			tags = append(tags, tag)
		}
	}
	return tags
}

var mentionRegex = regexp.MustCompile(`@(\w+)`)

func extractMentions(content string) []string {
	matches := mentionRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var mentions []string
	for _, m := range matches {
		username := m[1]
		if !seen[username] {
			seen[username] = true
			mentions = append(mentions, username)
		}
	}
	return mentions
}
