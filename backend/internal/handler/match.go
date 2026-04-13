package handler

import (
	"net/http"
	"strconv"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MatchHandler struct {
	matchService *service.MatchService
}

func NewMatchHandler(matchService *service.MatchService) *MatchHandler {
	return &MatchHandler{matchService: matchService}
}

func (h *MatchHandler) Create(c *gin.Context) {
	var req dto.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	match, err := h.matchService.CreateMatch(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(match))
}

func (h *MatchHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", "invalid match id"))
		return
	}

	match, err := h.matchService.GetMatch(c.Request.Context(), id)
	if err != nil || match == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", "match not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(match))
}

func (h *MatchHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	var playerID *uuid.UUID
	if pid := c.Query("player_id"); pid != "" {
		parsed, err := uuid.Parse(pid)
		if err == nil {
			playerID = &parsed
		}
	}

	matches, total, err := h.matchService.ListMatches(c.Request.Context(), playerID, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to list matches"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessWithMeta(matches, &dto.Meta{Page: page, PerPage: perPage, Total: total}))
}

func (h *MatchHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", "invalid match id"))
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	if err := h.matchService.DeleteMatch(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{"message": "match deleted"}))
}
