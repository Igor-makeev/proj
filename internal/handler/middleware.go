package handler

import (
	"net/http"
	"proj/internal/entities/myerrors"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

func (h *Handler) useridentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		myerrors.NewErorResponse(c, http.StatusUnauthorized, "error:empty auth header")
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		myerrors.NewErorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	userID, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		myerrors.NewErorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userID)

}
