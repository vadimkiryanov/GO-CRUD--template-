package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vadimkiryanov/GO-CRUD/pkg/service"
)

type Handler struct {
	services *service.Service
}

// Инициализация обработчиков
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (handler *Handler) testPing(ctx *gin.Context) {
	newErrorResponse(ctx, 400, "test")
}

// Инициализация роутеров
func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New() // создание роутера

	auth := router.Group("/test")
	{
		auth.POST("/second-level", h.testPing) // создание второго уровня роутера
	}

	return router
}
