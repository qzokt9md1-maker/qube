package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/kuzuokatakumi/qube/internal/config"
	"github.com/kuzuokatakumi/qube/internal/db"
	"github.com/kuzuokatakumi/qube/internal/handler"
	"github.com/kuzuokatakumi/qube/internal/middleware"
	"github.com/kuzuokatakumi/qube/internal/repository/postgres"
	"github.com/kuzuokatakumi/qube/internal/service"
	"github.com/kuzuokatakumi/qube/internal/ws"
)

func main() {
	cfg := config.Load()

	// Database
	pool, err := db.NewPostgresPool(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL")

	// Redis
	rdb, err := db.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("Connected to Redis")

	// WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Repositories
	userRepo := postgres.NewUserRepo(pool)
	postRepo := postgres.NewPostRepo(pool)
	followRepo := postgres.NewFollowRepo(pool)
	likeRepo := postgres.NewLikeRepo(pool)
	bookmarkRepo := postgres.NewBookmarkRepo(pool)
	convRepo := postgres.NewConversationRepo(pool)
	msgRepo := postgres.NewMessageRepo(pool)
	notifRepo := postgres.NewNotificationRepo(pool)
	sessionRepo := postgres.NewSessionRepo(pool)
	blockRepo := postgres.NewBlockRepo(pool)
	muteRepo := postgres.NewMuteRepo(pool)
	hashtagRepo := postgres.NewHashtagRepo(pool)
	cursorRepo := postgres.NewTimelineCursorRepo(pool)

	// Services
	notifService := service.NewNotificationService(notifRepo, hub)
	timelineService := service.NewTimelineService(rdb, postRepo, followRepo, cursorRepo)
	authService := service.NewAuthService(userRepo, sessionRepo, cfg.JWT)
	userService := service.NewUserService(userRepo, blockRepo, muteRepo)
	postService := service.NewPostService(postRepo, userRepo, hashtagRepo, likeRepo, bookmarkRepo, notifService, timelineService)
	followService := service.NewFollowService(followRepo, userRepo, blockRepo, notifService)
	dmService := service.NewDMService(convRepo, msgRepo, blockRepo, notifService, hub)

	// GraphQL Handler
	gqlHandler := &handler.GraphQLHandler{
		AuthService:     authService,
		UserService:     userService,
		PostService:     postService,
		FollowService:   followService,
		DMService:       dmService,
		NotifService:    notifService,
		TimelineService: timelineService,
	}

	// Router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Timeout(30 * time.Second))
	corsOrigins := []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8080", "https://qube.social"}
	if extra := os.Getenv("CORS_ORIGINS"); extra != "" {
		for _, o := range strings.Split(extra, ",") {
			corsOrigins = append(corsOrigins, strings.TrimSpace(o))
		}
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Auth(cfg.JWT.Secret))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// GraphQL endpoint
	r.Handle("/graphql", gqlHandler)

	// WebSocket endpoint
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		hub.HandleWebSocket(w, r, userID)
	})

	// Upload endpoint
	uploadDir := "./uploads"
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	uploadHandler := &handler.UploadHandler{
		UploadDir: uploadDir,
		BaseURL:   baseURL,
	}
	r.Handle("/upload", uploadHandler)

	// Serve uploaded files
	fileServer := http.FileServer(http.Dir(uploadDir))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	// Suppress unused variable warnings
	_ = hashtagRepo
	_ = cursorRepo

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Qube server starting on port %s (env: %s)", cfg.Server.Port, cfg.Server.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")

	// Ensure uuid is used
	_ = uuid.Nil
}
