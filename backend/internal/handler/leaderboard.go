package handler

import (
	"net/http"
	"strconv"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type LeaderboardHandler struct {
	userService *service.UserService
}

func NewLeaderboardHandler(userService *service.UserService) *LeaderboardHandler {
	return &LeaderboardHandler{userService: userService}
}

func (h *LeaderboardHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	users, total, err := h.userService.GetLeaderboard(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to load leaderboard"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMeta(users, &dto.Meta{Page: page, PerPage: perPage, Total: total}))
}
