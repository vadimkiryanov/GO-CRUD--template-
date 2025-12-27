package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
}

// NewServer создает новый сервер с Gin
func NewServer() *Server {
	// Настройка Gin
	router := gin.New()

	// Middleware по умолчанию
	router.Use(gin.Logger())   // ✅ Цветные логи запросов
	router.Use(gin.Recovery()) // ✅ Обработка паник
	port := viper.GetString("port")

	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        router,
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
		router: router,
	}
}

// Router возвращает *gin.Engine для регистрации маршрутов
func (s *Server) Router() *gin.Engine {
	return s.router
}

// Run запускает сервер
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown останавливает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
