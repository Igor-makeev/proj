package handler

import (
	"io"
	"net/http"
	"proj/internal/entities/myerrors"
	"proj/pkg/luhn"

	"github.com/gin-gonic/gin"
)

func (h *Handler) loadOrderNumber(c *gin.Context) {

	id, ok := c.Get(userCtx)
	if !ok && id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": myerrors.DontHaveAccess.Error()})
		return
	}

	number, err := io.ReadAll(c.Request.Body)

	if err != nil || string(number) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": myerrors.InvalidOrderInput.Error()})
		return
	}

	if !luhn.LuhnValidation(string(number)) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": myerrors.InvalidOrderNumber.Error()})
		return
	}

	if err := h.service.SaveOrder(c.Request.Context(), string(number), id.(int)); err != nil {

		switch err {
		case myerrors.ErrOrdUsrConfl:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case myerrors.ErrOrdOverLap:
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		}
		return

	}
	c.Status(http.StatusAccepted)
}

func (h *Handler) getOrdersList(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok && id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": myerrors.DontHaveAccess})
		return
	}

	ordersList, err := h.service.GetOrders(c.Request.Context(), id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if len(ordersList) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"Info": "Oredrs not found"})
		return
	}
	c.JSON(http.StatusOK, ordersList)

}

func (h *Handler) getBallance(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok && id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": myerrors.DontHaveAccess})
		return
	}
	accountState, err := h.service.GetBalance(c.Request.Context(), id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, *accountState)
}

func (h *Handler) withdrawRequest(c *gin.Context) {

}

func (h *Handler) getWithdrawInfo(c *gin.Context) {

}
