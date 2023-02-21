package handler

import (
	"io/ioutil"
	"net/http"
	"proj/pkg/luhn"

	"github.com/gin-gonic/gin"
)

func (h *Handler) loadOrerNumber(c *gin.Context) {
	number, err := ioutil.ReadAll(c.Request.Body)
	if err != nil || string(number) == "" {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	if ok := luhn.LuhnValidation(string(number)); !ok {
		c.String(http.StatusUnprocessableEntity, "Not valid number of order")
		return
	}
}

func (h *Handler) getOrdersList(c *gin.Context) {

}

func (h *Handler) getBallance(c *gin.Context) {

}

func (h *Handler) withdrawRequest(c *gin.Context) {

}

func (h *Handler) getWithdrawInfo(c *gin.Context) {

}
