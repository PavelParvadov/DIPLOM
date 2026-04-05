package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"happyhouse/backend/internal/config"
	"happyhouse/backend/internal/repository/postgres"
	"happyhouse/backend/internal/transport/http/handler"
	"happyhouse/backend/internal/usecase"
	"happyhouse/backend/pkg/auth"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping database: %v", err)
	}
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		log.Fatalf("create upload dir: %v", err)
	}

	tokenManager := auth.NewTokenManager(cfg.JWTSecret)
	userRepo := postgres.NewUserRepository(db)
	houseRepo := postgres.NewHouseRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	postRepo := postgres.NewPostRepository(db)
	commentRepo := postgres.NewCommentRepository(db)
	chatRepo := postgres.NewChatRepository(db)
	inviteRepo := postgres.NewInviteCodeRepository(db)
	refreshRepo := postgres.NewRefreshTokenRepository(db)

	authUC := usecase.NewAuthUseCase(userRepo, refreshRepo, tokenManager, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	houseUC := usecase.NewHouseUseCase(houseRepo, inviteRepo)
	categoryUC := usecase.NewCategoryUseCase(houseRepo, categoryRepo)
	postUC := usecase.NewPostUseCase(houseRepo, postRepo, categoryRepo)
	commentUC := usecase.NewCommentUseCase(houseRepo, postRepo, commentRepo)
	chatUC := usecase.NewChatUseCase(houseRepo, chatRepo)
	inviteUC := usecase.NewInviteCodeUseCase(houseRepo, inviteRepo)

	httpHandler := handler.New(authUC, houseUC, categoryUC, postUC, commentUC, chatUC, inviteUC, tokenManager, cfg.FrontendOrigin, cfg.UploadDir)
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      httpHandler.Router(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("happyhouse backend listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown server: %v", err)
	}
}
