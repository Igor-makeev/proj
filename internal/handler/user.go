package handler

import (
	"io"
	"net/http"
	"proj/internal/entities/models"
	"proj/internal/entities/myerrors"
	"proj/pkg/luhn"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) loadOrderNumber(c *gin.Context) {

	id, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": myerrors.ErrDontHaveAccess.Error()})
		return
	}

	number, err := io.ReadAll(c.Request.Body)

	if err != nil || len(number) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": myerrors.ErrInvalidOrderInput.Error()})
		return
	}
	logrus.Print(number)
	if !luhn.LuhnValidation(string(number)) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": myerrors.ErrInvalidOrderNumber.Error()})
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
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": myerrors.ErrDontHaveAccess.Error()})
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
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, ordersList)

}

func (h *Handler) getBallance(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": myerrors.ErrDontHaveAccess.Error()})
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
	id, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": myerrors.ErrDontHaveAccess.Error()})
		return
	}
	var withdraw models.Withdrawal
	if err := c.ShouldBindJSON(&withdraw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Withdraw(c.Request.Context(), withdraw, id.(int)); err != nil {

		switch err {

		case myerrors.ErrNoMoney:
			c.JSON(http.StatusPaymentRequired, err.Error())
		case myerrors.ErrInvalidOrderNumber:
			c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			c.JSON(http.StatusInternalServerError, err.Error())

		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"withdrawal": "done"})
}

func (h *Handler) getWithdrawInfo(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": myerrors.ErrDontHaveAccess.Error()})
		return
	}

	withdrawls, err := h.service.GetWithdrawals(c.Request.Context(), id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if len(withdrawls) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"info": "withdrawls not found"})
		return
	}
	c.JSON(http.StatusOK, withdrawls)
}
