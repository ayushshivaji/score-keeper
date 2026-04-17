package handler

import (
	"errors"
	"net/http"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QuickPairHandler struct {
	service *service.QuickPairService
}

func NewQuickPairHandler(s *service.QuickPairService) *QuickPairHandler {
	return &QuickPairHandler{service: s}
}

func (h *QuickPairHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	pairs, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to list quick pairs"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(pairs))
}

func (h *QuickPairHandler) Create(c *gin.Context) {
	var req dto.CreateQuickPairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", "player1_id and player2_id are required"))
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	qp, err := h.service.Create(c.Request.Context(), userID, req.Player1ID, req.Player2ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSamePlayer):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", err.Error()))
		case errors.Is(err, service.ErrPlayerNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", err.Error()))
		case errors.Is(err, service.ErrQuickPairExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse("CONFLICT", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to create quick pair"))
		}
		return
	}
	c.JSON(http.StatusCreated, dto.Success(qp))
}

func (h *QuickPairHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("BAD_REQUEST", "invalid quick pair id"))
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	if err := h.service.Delete(c.Request.Context(), userID, id); err != nil {
		if errors.Is(err, service.ErrQuickPairNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("NOT_FOUND", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "failed to delete quick pair"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(gin.H{"message": "quick pair deleted"}))
}
