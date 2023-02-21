package handler

import (
	"io/ioutil"
	"net/http"
	"proj/internal/entities/myerrors"
	"proj/pkg/luhn"

	"github.com/gin-gonic/gin"
)

func (h *Handler) loadOrderNumber(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok && id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": myerrors.DontHaveAccess})
		return
	}

	number, err := ioutil.ReadAll(c.Request.Body)

	if err != nil || string(number) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if ok := luhn.LuhnValidation(string(number)); !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": myerrors.InvalidOrderNumber})
		return
	}

	if err := h.service.SaveOrder(c.Request.Context(), string(number), id.(int)); err != nil {

		switch err {
		case myerrors.ErrOrdUsrConfl:
			c.JSON(http.StatusConflict, gin.H{"error": err})
		case myerrors.ErrOrdOverLap:
			c.JSON(http.StatusOK, gin.H{"error": err})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})

		}
		return

	}
	c.Status(http.StatusAccepted)
}

func (h *Handler) getOrdersList(c *gin.Context) {

}

func (h *Handler) getBallance(c *gin.Context) {

}

func (h *Handler) withdrawRequest(c *gin.Context) {

}

func (h *Handler) getWithdrawInfo(c *gin.Context) {

}
