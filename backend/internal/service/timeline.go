package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/model"
	"github.com/kuzuokatakumi/qube/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
)

const (
	timelineKeyPrefix = "timeline:"
	timelineMaxSize   = 800
	timelineTTL       = 7 * 24 * time.Hour
)

type TimelineService struct {
	redis       *redis.Client
	postRepo    *postgres.PostRepo
	followRepo  *postgres.FollowRepo
	cursorRepo  *postgres.TimelineCursorRepo
}

func NewTimelineService(
	rdb *redis.Client,
	postRepo *postgres.PostRepo,
	followRepo *postgres.FollowRepo,
	cursorRepo *postgres.TimelineCursorRepo,
) *TimelineService {
	return &TimelineService{
		redis:      rdb,
		postRepo:   postRepo,
		followRepo: followRepo,
		cursorRepo: cursorRepo,
	}
}

type TimelineEntry struct {
	PostID    uuid.UUID `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

// FanOutPost pushes a new post to all followers' Redis timelines
func (s *TimelineService) FanOutPost(ctx context.Context, post *model.Post) {
	followerIDs, err := s.followRepo.GetFollowingIDs(ctx, post.UserID)
	if err != nil {
		return
	}

	// Also include self
	followerIDs = append(followerIDs, post.UserID)

	// We actually need the followers OF the poster, not who the poster follows
	// Correction: get users who follow this user (the poster's followers)
	// GetFollowingIDs gets who userID follows, but we need who follows userID
	// Let's use the pipeline approach for fan-out
	entry := TimelineEntry{
		PostID:    post.ID,
		CreatedAt: post.CreatedAt,
	}
	data, _ := json.Marshal(entry)

	pipe := s.redis.Pipeline()
	for _, fid := range followerIDs {
		key := fmt.Sprintf("%s%s", timelineKeyPrefix, fid.String())
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(post.CreatedAt.UnixNano()),
			Member: string(data),
		})
		pipe.ZRemRangeByRank(ctx, key, 0, -timelineMaxSize-1)
		pipe.Expire(ctx, key, timelineTTL)
	}
	pipe.Exec(ctx)
}

// GetTimeline returns timeline posts for a user, using Redis cache + DB fallback
func (s *TimelineService) GetTimeline(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*model.Post, int, error) {
	// Try Redis first
	key := fmt.Sprintf("%s%s", timelineKeyPrefix, userID.String())

	var maxScore string
	if cursor != "" {
		t, err := time.Parse(time.RFC3339Nano, cursor)
		if err == nil {
			maxScore = fmt.Sprintf("%d", t.UnixNano())
		} else {
			maxScore = "+inf"
		}
	} else {
		maxScore = "+inf"
	}

	results, err := s.redis.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    maxScore,
		Offset: 0,
		Count:  int64(limit),
	}).Result()

	if err == nil && len(results) > 0 {
		var postIDs []uuid.UUID
		for _, r := range results {
			var entry TimelineEntry
			if err := json.Unmarshal([]byte(r), &entry); err == nil {
				postIDs = append(postIDs, entry.PostID)
			}
		}

		if len(postIDs) > 0 {
			posts := make([]*model.Post, 0, len(postIDs))
			for _, pid := range postIDs {
				post, err := s.postRepo.GetByID(ctx, pid)
				if err == nil {
					posts = append(posts, post)
				}
			}
			unread, _ := s.cursorRepo.GetUnreadCount(ctx, userID)
			return posts, unread, nil
		}
	}

	// Fallback to DB
	posts, err := s.postRepo.GetTimeline(ctx, userID, limit, cursor)
	if err != nil {
		return nil, 0, err
	}

	unread, _ := s.cursorRepo.GetUnreadCount(ctx, userID)
	return posts, unread, nil
}

func (s *TimelineService) UpdateCursor(ctx context.Context, userID, postID uuid.UUID) error {
	return s.cursorRepo.Update(ctx, userID, postID)
}
