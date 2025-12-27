package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vadimkiryanov/GO-CRUD/internal/service"
)

type Handler struct {
	services *service.Service
}

// Инициализация обработчиков
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// Инициализация роутеров

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New() // создание роутера

	router.Use(gin.Logger())   // ✅ Цветные логи запросов
	router.Use(gin.Recovery()) // ✅ Обработка паник

	// ❌ ВРЕМЕННО ДЛЯ DEV (удалите в проде)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp) // регистрация
		auth.POST("/sign-in", h.signIn) // авторизация
	}

	return router
}
