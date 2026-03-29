package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kuzuokatakumi/qube/internal/service"
)

type GraphQLHandler struct {
	AuthService   *service.AuthService
	UserService   *service.UserService
	PostService   *service.PostService
	FollowService *service.FollowService
	DMService     *service.DMService
	NotifService  *service.NotificationService
	TimelineService *service.TimelineService
}

type graphqlRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type graphqlResponse struct {
	Data   interface{}    `json:"data,omitempty"`
	Errors []graphqlError `json:"errors,omitempty"`
}

type graphqlError struct {
	Message string `json:"message"`
}

func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req graphqlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result := h.executeQuery(r, req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(graphqlResponse{
		Errors: []graphqlError{{Message: msg}},
	})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
