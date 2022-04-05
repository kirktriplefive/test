package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kirktriplefive/test/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler  {
	return &Handler{services: services}
}



func (h *Handler) InitRoutes() *gin.Engine {
	router:=gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api")
	{
		order:=api.Group("/order")
		{
			order.POST("/", h.createOrder)
			order.GET("/:id", h.getOrderById)
			items:=order.Group(":id/items")
			{
				items.POST("/", h.createOrderWithItem)
				items.POST("/new", h.createOrderWithNewItem)
			}
		}
		item:=api.Group("/item")
		{
			item.POST("/", h.createItem)
			item.GET("/", h.getAllItem)
		}
	}
	return router
}