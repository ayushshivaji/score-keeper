package handler

import (
	"net/http"
	"strconv"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) List(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	users, total, err := h.userService.ListUsers(c.Request.Context(), search, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to list users"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMeta(users, &dto.Meta{Page: page, PerPage: perPage, Total: total}))
}

func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", "invalid user id"))
		return
	}

	profile, err := h.userService.GetUserProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to load profile"))
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", "user not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(profile))
}

func (h *UserHandler) HeadToHead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", "invalid user id"))
		return
	}
	opponentID, err := uuid.Parse(c.Param("opponentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", "invalid opponent id"))
		return
	}

	h2h, err := h.userService.GetHeadToHead(c.Request.Context(), id, opponentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(h2h))
}
