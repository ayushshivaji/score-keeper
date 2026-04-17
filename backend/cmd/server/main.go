package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ayush-sr/score-keeper/backend/internal/config"
	"github.com/ayush-sr/score-keeper/backend/internal/handler"
	"github.com/ayush-sr/score-keeper/backend/internal/repository"
	"github.com/ayush-sr/score-keeper/backend/internal/router"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load .env in development
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Database
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}
	log.Info().Msg("connected to database")

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	matchRepo := repository.NewMatchRepository(pool)
	quickPairRepo := repository.NewQuickPairRepository(pool)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo, matchRepo)
	matchService := service.NewMatchService(matchRepo, userRepo)
	quickPairService := service.NewQuickPairService(quickPairRepo, userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, cfg)
	userHandler := handler.NewUserHandler(userService)
	matchHandler := handler.NewMatchHandler(matchService)
	leaderboardHandler := handler.NewLeaderboardHandler(userService)
	quickPairHandler := handler.NewQuickPairHandler(quickPairService)

	// Router
	r := gin.Default()
	router.Setup(r, authHandler, userHandler, matchHandler, leaderboardHandler, quickPairHandler, authService, cfg.FrontendURL)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Info().Str("addr", addr).Msg("starting server")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
