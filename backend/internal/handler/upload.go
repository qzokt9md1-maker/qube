package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/middleware"
)

type UploadHandler struct {
	UploadDir string
	BaseURL   string
}

type uploadResponse struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	MediaType    string `json:"mediaType"`
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Max 10MB
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "no file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	mediaType := ""
	switch {
	case strings.HasPrefix(contentType, "image/"):
		mediaType = "image"
	case strings.HasPrefix(contentType, "video/"):
		mediaType = "video"
	case contentType == "image/gif":
		mediaType = "gif"
	default:
		http.Error(w, "unsupported file type", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102"), hex.EncodeToString(randBytes), ext)

	// Ensure upload directory exists
	if err := os.MkdirAll(h.UploadDir, 0755); err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	// Save file
	destPath := filepath.Join(h.UploadDir, filename)
	dest, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	url := h.BaseURL + "/uploads/" + filename

	resp := uploadResponse{
		ID:        uuid.New().String(),
		URL:       url,
		MediaType: mediaType,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
