package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ayush-sr/score-keeper/backend/internal/config"
	"github.com/ayush-sr/score-keeper/backend/internal/dto"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	authService *service.AuthService
	oauthConfig *oauth2.Config
	cfg         *config.Config
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

type googleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", "missing code"))
		return
	}

	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("OAUTH_ERROR", "failed to exchange token"))
		return
	}

	client := h.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("OAUTH_ERROR", "failed to get user info"))
		return
	}
	defer resp.Body.Close()

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("OAUTH_ERROR", "failed to decode user info"))
		return
	}

	var avatarURL *string
	if userInfo.Picture != "" {
		avatarURL = &userInfo.Picture
	}

	user, err := h.authService.UpsertUser(c.Request.Context(), userInfo.Sub, userInfo.Email, userInfo.Name, avatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to create user"))
		return
	}

	accessToken, err := h.authService.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to generate access token"))
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to generate refresh token"))
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, h.cfg.FrontendURL+"/dashboard")
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", "missing refresh token"))
		return
	}

	accessToken, newRefresh, err := h.authService.RefreshAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("UNAUTHORIZED", err.Error()))
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", newRefresh, 604800, "/", "", false, true)
	c.JSON(http.StatusOK, dto.Success(gin.H{"message": "tokens refreshed"}))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	_ = h.authService.Logout(c.Request.Context(), userID)

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, dto.Success(gin.H{"message": "logged out"}))
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	user, err := h.authService.GetUser(c.Request.Context(), userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", "user not found"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(user))
}
