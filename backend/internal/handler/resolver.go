package handler

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/middleware"
)

// executeQuery is a simple GraphQL executor that parses operation names
// and routes to the appropriate resolver. For production, replace with gqlgen.
func (h *GraphQLHandler) executeQuery(r *http.Request, req graphqlRequest) graphqlResponse {
	ctx := r.Context()
	op := strings.ToLower(req.OperationName)

	// Helper to get authenticated user
	getUserID := func() (uuid.UUID, *graphqlResponse) {
		uid, err := middleware.RequireAuth(ctx)
		if err != nil {
			return uuid.Nil, &graphqlResponse{Errors: []graphqlError{{Message: "unauthorized"}}}
		}
		return uid, nil
	}

	// Simple operation routing based on operationName
	switch op {

	// ==================== Auth ====================
	case "register":
		vars := req.Variables
		input, _ := vars["input"].(map[string]interface{})
		username, _ := input["username"].(string)
		displayName, _ := input["displayName"].(string)
		email, _ := input["email"].(string)
		password, _ := input["password"].(string)

		payload, err := h.AuthService.Register(ctx, username, displayName, email, password)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{
			"register": map[string]interface{}{
				"accessToken":  payload.AccessToken,
				"refreshToken": payload.RefreshToken,
				"user":         userToMap(payload.User),
			},
		}}

	case "login":
		vars := req.Variables
		input, _ := vars["input"].(map[string]interface{})
		email, _ := input["email"].(string)
		password, _ := input["password"].(string)

		payload, err := h.AuthService.Login(ctx, email, password)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{
			"login": map[string]interface{}{
				"accessToken":  payload.AccessToken,
				"refreshToken": payload.RefreshToken,
				"user":         userToMap(payload.User),
			},
		}}

	case "refreshtoken":
		token, _ := req.Variables["token"].(string)
		payload, err := h.AuthService.RefreshToken(ctx, token)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{
			"refreshToken": map[string]interface{}{
				"accessToken":  payload.AccessToken,
				"refreshToken": payload.RefreshToken,
				"user":         userToMap(payload.User),
			},
		}}

	// ==================== User ====================
	case "me":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		user, err := h.UserService.GetByID(ctx, uid)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"me": userToMap(user)}}

	case "user":
		username, _ := req.Variables["username"].(string)
		user, err := h.UserService.GetByUsername(ctx, username)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: "user not found"}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"user": userToMap(user)}}

	case "updateprofile":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		input, _ := req.Variables["input"].(map[string]interface{})
		var displayName, bio, location, website *string
		var isPrivate *bool
		if v, ok := input["displayName"].(string); ok {
			displayName = &v
		}
		if v, ok := input["bio"].(string); ok {
			bio = &v
		}
		if v, ok := input["location"].(string); ok {
			location = &v
		}
		if v, ok := input["website"].(string); ok {
			website = &v
		}
		if v, ok := input["isPrivate"].(bool); ok {
			isPrivate = &v
		}
		user, err := h.UserService.UpdateProfile(ctx, uid, displayName, bio, location, website, isPrivate)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"updateProfile": userToMap(user)}}

	case "searchusers":
		query, _ := req.Variables["query"].(string)
		limit := intFromVar(req.Variables, "limit", 20)
		cursor, _ := req.Variables["cursor"].(string)
		users, err := h.UserService.Search(ctx, query, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		userMaps := make([]interface{}, len(users))
		for i, u := range users {
			userMaps[i] = userToMap(u)
		}
		return graphqlResponse{Data: map[string]interface{}{"searchUsers": map[string]interface{}{
			"users":   userMaps,
			"hasMore": len(users) == limit,
		}}}

	// ==================== Posts ====================
	case "createpost":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		input, _ := req.Variables["input"].(map[string]interface{})
		content, _ := input["content"].(string)
		var replyToID, quoteOfID *uuid.UUID
		if v, ok := input["replyToId"].(string); ok {
			id, _ := uuid.Parse(v)
			replyToID = &id
		}
		if v, ok := input["quoteOfId"].(string); ok {
			id, _ := uuid.Parse(v)
			quoteOfID = &id
		}
		post, err := h.PostService.Create(ctx, uid, content, replyToID, quoteOfID, nil)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"createPost": postToMap(post)}}

	case "deletepost":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		postID, _ := uuid.Parse(req.Variables["id"].(string))
		if err := h.PostService.Delete(ctx, uid, postID); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"deletePost": true}}

	case "post":
		postID, _ := uuid.Parse(req.Variables["id"].(string))
		post, err := h.PostService.GetByID(ctx, postID)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: "post not found"}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"post": postToMap(post)}}

	case "timeline":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		limit := intFromVar(req.Variables, "limit", 20)
		cursor, _ := req.Variables["cursor"].(string)
		posts, unread, err := h.TimelineService.GetTimeline(ctx, uid, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		postMaps := make([]interface{}, len(posts))
		for i, p := range posts {
			postMaps[i] = postToMap(p)
		}
		var nextCursor string
		if len(posts) > 0 {
			nextCursor = posts[len(posts)-1].CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00")
		}
		return graphqlResponse{Data: map[string]interface{}{"timeline": map[string]interface{}{
			"posts":       postMaps,
			"hasMore":     len(posts) == limit,
			"cursor":      nextCursor,
			"unreadCount": unread,
		}}}

	case "userposts":
		username, _ := req.Variables["username"].(string)
		limit := intFromVar(req.Variables, "limit", 20)
		cursor, _ := req.Variables["cursor"].(string)
		user, err := h.UserService.GetByUsername(ctx, username)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: "user not found"}}}
		}
		posts, err := h.PostService.GetUserPosts(ctx, user.ID, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		postMaps := make([]interface{}, len(posts))
		for i, p := range posts {
			postMaps[i] = postToMap(p)
		}
		return graphqlResponse{Data: map[string]interface{}{"userPosts": map[string]interface{}{
			"posts":   postMaps,
			"hasMore": len(posts) == limit,
		}}}

	case "postreplies":
		postID, _ := uuid.Parse(req.Variables["postId"].(string))
		limit := intFromVar(req.Variables, "limit", 20)
		cursor, _ := req.Variables["cursor"].(string)
		posts, err := h.PostService.GetReplies(ctx, postID, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		postMaps := make([]interface{}, len(posts))
		for i, p := range posts {
			postMaps[i] = postToMap(p)
		}
		return graphqlResponse{Data: map[string]interface{}{"postReplies": map[string]interface{}{
			"posts":   postMaps,
			"hasMore": len(posts) == limit,
		}}}

	// ==================== Likes ====================
	case "likepost":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		postID, _ := uuid.Parse(req.Variables["postId"].(string))
		post, err := h.PostService.Like(ctx, uid, postID)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"likePost": postToMap(post)}}

	case "unlikepost":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		postID, _ := uuid.Parse(req.Variables["postId"].(string))
		post, err := h.PostService.Unlike(ctx, uid, postID)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"unlikePost": postToMap(post)}}

	case "repost":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		postID, _ := uuid.Parse(req.Variables["postId"].(string))
		post, err := h.PostService.Repost(ctx, uid, postID)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"repost": postToMap(post)}}

	case "bookmarkpost":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		postID, _ := uuid.Parse(req.Variables["postId"].(string))
		if err := h.PostService.Bookmark(ctx, uid, postID); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"bookmarkPost": true}}

	case "unbookmarkpost":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		postID, _ := uuid.Parse(req.Variables["postId"].(string))
		if err := h.PostService.Unbookmark(ctx, uid, postID); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"unbookmarkPost": true}}

	case "bookmarks":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		limit := intFromVar(req.Variables, "limit", 20)
		cursor, _ := req.Variables["cursor"].(string)
		posts, err := h.PostService.GetBookmarks(ctx, uid, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		postMaps := make([]interface{}, len(posts))
		for i, p := range posts {
			postMaps[i] = postToMap(p)
		}
		return graphqlResponse{Data: map[string]interface{}{"bookmarks": map[string]interface{}{
			"posts":   postMaps,
			"hasMore": len(posts) == limit,
		}}}

	// ==================== Follow ====================
	case "follow":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		targetID, _ := uuid.Parse(req.Variables["userId"].(string))
		user, err := h.FollowService.Follow(ctx, uid, targetID)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"follow": userToMap(user)}}

	case "unfollow":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		targetID, _ := uuid.Parse(req.Variables["userId"].(string))
		user, err := h.FollowService.Unfollow(ctx, uid, targetID)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"unfollow": userToMap(user)}}

	case "followers":
		username, _ := req.Variables["username"].(string)
		limit := intFromVar(req.Variables, "limit", 20)
		cursor, _ := req.Variables["cursor"].(string)
		user, err := h.UserService.GetByUsername(ctx, username)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: "user not found"}}}
		}
		users, err := h.FollowService.GetFollowers(ctx, user.ID, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		userMaps := make([]interface{}, len(users))
		for i, u := range users {
			userMaps[i] = userToMap(u)
		}
		return graphqlResponse{Data: map[string]interface{}{"followers": map[string]interface{}{
			"users":   userMaps,
			"hasMore": len(users) == limit,
		}}}

	case "following":
		username, _ := req.Variables["username"].(string)
		limit := intFromVar(req.Variables, "limit", 20)
		cursor, _ := req.Variables["cursor"].(string)
		user, err := h.UserService.GetByUsername(ctx, username)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: "user not found"}}}
		}
		users, err := h.FollowService.GetFollowing(ctx, user.ID, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		userMaps := make([]interface{}, len(users))
		for i, u := range users {
			userMaps[i] = userToMap(u)
		}
		return graphqlResponse{Data: map[string]interface{}{"following": map[string]interface{}{
			"users":   userMaps,
			"hasMore": len(users) == limit,
		}}}

	// ==================== DM ====================
	case "createconversation":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		input, _ := req.Variables["input"].(map[string]interface{})
		message, _ := input["message"].(string)
		pIDs := parseUUIDList(input["participantIds"])
		conv, err := h.DMService.CreateConversation(ctx, uid, pIDs, message)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"createConversation": convToMap(conv)}}

	case "sendmessage":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		input, _ := req.Variables["input"].(map[string]interface{})
		convID, _ := uuid.Parse(input["conversationId"].(string))
		content, _ := input["content"].(string)
		msg, err := h.DMService.SendMessage(ctx, uid, convID, content)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"sendMessage": msgToMap(msg)}}

	case "conversations":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		limit := intFromVar(req.Variables, "limit", 20)
		cursor, _ := req.Variables["cursor"].(string)
		convs, err := h.DMService.GetConversations(ctx, uid, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		convMaps := make([]interface{}, len(convs))
		for i, c := range convs {
			convMaps[i] = convToMap(c)
		}
		return graphqlResponse{Data: map[string]interface{}{"conversations": map[string]interface{}{
			"conversations": convMaps,
			"hasMore":       len(convs) == limit,
		}}}

	case "messages":
		convID, _ := uuid.Parse(req.Variables["conversationId"].(string))
		limit := intFromVar(req.Variables, "limit", 30)
		cursor, _ := req.Variables["cursor"].(string)
		msgs, err := h.DMService.GetMessages(ctx, convID, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		msgMaps := make([]interface{}, len(msgs))
		for i, m := range msgs {
			msgMaps[i] = msgToMap(m)
		}
		return graphqlResponse{Data: map[string]interface{}{"messages": map[string]interface{}{
			"messages": msgMaps,
			"hasMore":  len(msgs) == limit,
		}}}

	case "markconversationread":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		convID, _ := uuid.Parse(req.Variables["conversationId"].(string))
		if err := h.DMService.MarkRead(ctx, convID, uid); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"markConversationRead": true}}

	// ==================== Notifications ====================
	case "notifications":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		limit := intFromVar(req.Variables, "limit", 20)
		cursor, _ := req.Variables["cursor"].(string)
		notifs, unread, err := h.NotifService.GetByUserID(ctx, uid, limit, cursor)
		if err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		notifMaps := make([]interface{}, len(notifs))
		for i, n := range notifs {
			notifMaps[i] = notifToMap(n)
		}
		return graphqlResponse{Data: map[string]interface{}{"notifications": map[string]interface{}{
			"notifications": notifMaps,
			"unreadCount":   unread,
			"hasMore":       len(notifs) == limit,
		}}}

	case "marknotificationsread":
		ids := parseUUIDList(req.Variables["ids"])
		if err := h.NotifService.MarkRead(ctx, ids); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"markNotificationsRead": true}}

	case "markallnotificationsread":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		if err := h.NotifService.MarkAllRead(ctx, uid); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"markAllNotificationsRead": true}}

	// ==================== Block / Mute ====================
	case "blockuser":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		targetID, _ := uuid.Parse(req.Variables["userId"].(string))
		if err := h.UserService.Block(ctx, uid, targetID); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"blockUser": true}}

	case "unblockuser":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		targetID, _ := uuid.Parse(req.Variables["userId"].(string))
		if err := h.UserService.Unblock(ctx, uid, targetID); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"unblockUser": true}}

	case "muteuser":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		targetID, _ := uuid.Parse(req.Variables["userId"].(string))
		if err := h.UserService.Mute(ctx, uid, targetID); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"muteUser": true}}

	case "unmuteuser":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		targetID, _ := uuid.Parse(req.Variables["userId"].(string))
		if err := h.UserService.Unmute(ctx, uid, targetID); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"unmuteUser": true}}

	// ==================== Timeline Cursor ====================
	case "updatetimelinecursor":
		uid, errResp := getUserID()
		if errResp != nil {
			return *errResp
		}
		postID, _ := uuid.Parse(req.Variables["lastSeenPostId"].(string))
		if err := h.TimelineService.UpdateCursor(ctx, uid, postID); err != nil {
			return graphqlResponse{Errors: []graphqlError{{Message: err.Error()}}}
		}
		return graphqlResponse{Data: map[string]interface{}{"updateTimelineCursor": true}}

	default:
		return graphqlResponse{Errors: []graphqlError{{Message: "unknown operation: " + req.OperationName}}}
	}
}

func intFromVar(vars map[string]interface{}, key string, fallback int) int {
	if v, ok := vars[key].(float64); ok {
		return int(v)
	}
	return fallback
}

func parseUUIDList(v interface{}) []uuid.UUID {
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	ids := make([]uuid.UUID, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			if id, err := uuid.Parse(s); err == nil {
				ids = append(ids, id)
			}
		}
	}
	return ids
}
