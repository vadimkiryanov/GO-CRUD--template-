package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/vadimkiryanov/GO-CRUD/internal/core/server"
	"github.com/vadimkiryanov/GO-CRUD/internal/features/auth"
)

type App struct {
	server *server.Server
}

func New(db *sqlx.DB) (*App, error) {
	// Создание сервера
	srv := server.NewServer()

	// Инициализация зависимостей auth
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// Регистрация маршрутов
	authHandler.InitRouters(srv.Router())

	return &App{
		server: srv,
	}, nil
}

func (a *App) Run() error {
	return a.server.Run()
}
