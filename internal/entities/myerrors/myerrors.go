package myerrors

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	InvalidLoginOrPassword = errors.New("error: invalid login or password")
	DontHaveAccess         = errors.New("error: dont have access")
	InvalidOrderNumber     = errors.New("error: invalid order number")
	ErrOrdOverLap          = errors.New("error: order already exist")
	ErrOrdUsrConfl         = errors.New("error: order was added by other customer")
)

type LoginConflict struct {
	Elem string
}

func (lc *LoginConflict) Error() string {
	return fmt.Sprintf("error: user with login:%v, has already exists", lc.Elem)
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
