package api

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

// Запуск сервера
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,          // 1 MB
		ReadTimeout:    10 * time.Second, // 10 seconds
		WriteTimeout:   10 * time.Second, // 10 seconds
	}

	return s.httpServer.ListenAndServe() // Бесконечный цикл для прослушивания всех запросов и последующей обработки
}

// Остановка сервера
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) // Остановка сервера
}
