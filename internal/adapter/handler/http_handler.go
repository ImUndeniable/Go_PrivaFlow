// WAITER
package handler

import (
	"net/http"

	"github.com/ImUndeniable/Go_PrivaFlow/internal/core/services"
	"github.com/gin-gonic/gin"
)

type ErasureHandler struct {
	svc *services.ErasureService
}

func NewErasureHandler(svc *services.ErasureService) *ErasureHandler {
	return &ErasureHandler{svc: svc}
}

func (h *ErasureHandler) RequestErasure(c *gin.Context) {
	// 1. Waiter writes down the order (JSON -> Struct)
	var input struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Waiter YELLS to the Chef (The "Missing" Link!)
	// 👇 THIS IS THE LINE YOU WERE LOOKING FOR 👇

	req, err := h.svc.RequestErasure(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": req})
}
