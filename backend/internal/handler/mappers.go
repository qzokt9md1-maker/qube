package handler

import (
	"github.com/kuzuokatakumi/qube/internal/model"
)

func userToMap(u *model.User) map[string]interface{} {
	if u == nil {
		return nil
	}
	return map[string]interface{}{
		"id":             u.ID.String(),
		"username":       u.Username,
		"displayName":    u.DisplayName,
		"bio":            u.Bio,
		"avatarUrl":      u.AvatarURL,
		"headerUrl":      u.HeaderURL,
		"location":       u.Location,
		"website":        u.Website,
		"isVerified":     u.IsVerified,
		"isPrivate":      u.IsPrivate,
		"followerCount":  u.FollowerCount,
		"followingCount": u.FollowingCount,
		"postCount":      u.PostCount,
		"createdAt":      u.CreatedAt,
	}
}

func postToMap(p *model.Post) map[string]interface{} {
	if p == nil {
		return nil
	}
	m := map[string]interface{}{
		"id":          p.ID.String(),
		"content":     p.Content,
		"likeCount":   p.LikeCount,
		"repostCount": p.RepostCount,
		"replyCount":  p.ReplyCount,
		"quoteCount":  p.QuoteCount,
		"createdAt":   p.CreatedAt,
		"user":        userToMap(p.User),
		"media":       mediaListToMap(p.Media),
	}
	if p.ReplyToID != nil {
		m["replyToId"] = p.ReplyToID.String()
	}
	if p.RepostOfID != nil {
		m["repostOfId"] = p.RepostOfID.String()
	}
	if p.QuoteOfID != nil {
		m["quoteOfId"] = p.QuoteOfID.String()
	}
	return m
}

func mediaListToMap(media []model.Media) []map[string]interface{} {
	result := make([]map[string]interface{}, len(media))
	for i, m := range media {
		result[i] = map[string]interface{}{
			"id":           m.ID.String(),
			"mediaType":    m.MediaType,
			"url":          m.URL,
			"thumbnailUrl": m.ThumbnailURL,
			"width":        m.Width,
			"height":       m.Height,
		}
	}
	return result
}

func convToMap(c *model.Conversation) map[string]interface{} {
	if c == nil {
		return nil
	}
	participants := make([]interface{}, len(c.Participants))
	for i, p := range c.Participants {
		participants[i] = userToMap(&p)
	}
	result := map[string]interface{}{
		"id":           c.ID.String(),
		"isGroup":      c.IsGroup,
		"name":         c.Name,
		"participants": participants,
		"updatedAt":    c.UpdatedAt,
	}
	if c.LastMessage != nil {
		result["lastMessage"] = msgToMap(c.LastMessage)
	}
	return result
}

func msgToMap(m *model.Message) map[string]interface{} {
	if m == nil {
		return nil
	}
	return map[string]interface{}{
		"id":             m.ID.String(),
		"conversationId": m.ConversationID.String(),
		"sender":         userToMap(m.Sender),
		"content":        m.Content,
		"createdAt":      m.CreatedAt,
	}
}

func notifToMap(n *model.Notification) map[string]interface{} {
	if n == nil {
		return nil
	}
	result := map[string]interface{}{
		"id":        n.ID.String(),
		"actor":     userToMap(n.Actor),
		"type":      n.Type,
		"isRead":    n.IsRead,
		"createdAt": n.CreatedAt,
	}
	if n.PostID != nil {
		result["postId"] = n.PostID.String()
	}
	return result
}
