package myerrors

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LoginConflict struct {
	Elem string
}

func (lc *LoginConflict) Error() string {
	return fmt.Sprintf("error: user with login:%v, has already exists", lc.Elem)
}

type InvalidLoginOrPassword struct {
}

func (ilop *InvalidLoginOrPassword) Error() string {
	return "error: wrong login or password"
}

type statusResponse struct {
	Status string `json:"status"`
}
type errorResponse struct {
	Message string `json:"message"`
}

func NewErorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Errorf(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
