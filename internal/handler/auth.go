package handler

import (
	"errors"
	"fmt"
	"net/http"
	"proj/internal/entities/models"
	"proj/internal/entities/myerrors"

	"github.com/gin-gonic/gin"
)

func (h *Handler) register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Authorization.CreateUser(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(*myerrors.LoginConflict); ok {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.Authorization.GenerateToken(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		if ok := errors.Is(myerrors.ErrInvalidLoginOrPassword, err); ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header(authorizationHeader, fmt.Sprintf("Bearer %s", token))
	c.Status(http.StatusOK)
}

func (h *Handler) login(c *gin.Context) {
	var input models.SignInInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Authorization.GenerateToken(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		if ok := errors.Is(myerrors.ErrInvalidLoginOrPassword, err); ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header(authorizationHeader, fmt.Sprintf("Bearer %s", token))
	c.Status(http.StatusOK)
}
