package router

import (
	"github.com/ayush-sr/score-keeper/backend/internal/handler"
	"github.com/ayush-sr/score-keeper/backend/internal/middleware"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func Setup(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	matchHandler *handler.MatchHandler,
	leaderboardHandler *handler.LeaderboardHandler,
	authService *service.AuthService,
	frontendURL string,
) {
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS(frontendURL))

	v1 := r.Group("/api/v1")

	// Public auth routes
	auth := v1.Group("/auth")
	auth.GET("/google", authHandler.GoogleLogin)
	auth.GET("/google/callback", authHandler.GoogleCallback)
	auth.POST("/login", authHandler.StaticLogin)
	auth.POST("/refresh", authHandler.Refresh)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthRequired(authService))

	protected.POST("/auth/logout", authHandler.Logout)
	protected.GET("/auth/me", authHandler.Me)

	protected.GET("/users", userHandler.List)
	protected.GET("/users/:id", userHandler.Get)
	protected.GET("/users/:id/head-to-head/:opponentId", userHandler.HeadToHead)

	protected.GET("/leaderboard", leaderboardHandler.List)

	protected.POST("/matches", matchHandler.Create)
	protected.GET("/matches", matchHandler.List)
	protected.GET("/matches/:id", matchHandler.Get)
	protected.DELETE("/matches/:id", matchHandler.Delete)
}
