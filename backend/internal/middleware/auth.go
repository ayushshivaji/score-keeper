package middleware

import (
	"net/http"
	"strings"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func AuthRequired(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		// Check cookie first
		if cookie, err := c.Cookie("access_token"); err == nil {
			tokenString = cookie
		}

		// Fall back to Authorization header
		if tokenString == "" {
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tokenString = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "missing access token"))
			c.Abort()
			return
		}

		claims, err := authService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "invalid or expired token"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
