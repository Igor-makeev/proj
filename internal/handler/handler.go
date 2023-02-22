package handler

import (
	"proj/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
	Router  *gin.Engine
}

func NewHandler(service *service.Service) *Handler {
	handler := &Handler{
		Router:  gin.New(),
		service: service,
	}

	api := handler.Router.Group("/api")
	{
		auth := api.Group("/user")
		{
			auth.POST("/register", handler.register)
			auth.POST("/login", handler.login)
		}

		user := api.Group("/user").Use(handler.useridentity)
		{
			user.POST("/orders", handler.loadOrderNumber)
			user.GET("/orders", handler.getOrdersList)
			user.GET("/balance", handler.getBallance)
			user.POST("/balance/withdraw", handler.withdrawRequest)
			user.GET("/withdrawals", handler.getWithdrawInfo)

		}

	}

	return handler
}
