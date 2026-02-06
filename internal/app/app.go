package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vadimkiryanov/GO-CRUD/internal/core/db"
	"github.com/vadimkiryanov/GO-CRUD/internal/core/handlers"
	"github.com/vadimkiryanov/GO-CRUD/internal/core/server"
	"github.com/vadimkiryanov/GO-CRUD/internal/features/auth"

	postsR "github.com/vadimkiryanov/GO-CRUD/internal/features/posts/repository"
	postsS "github.com/vadimkiryanov/GO-CRUD/internal/features/posts/service"
	postsT "github.com/vadimkiryanov/GO-CRUD/internal/features/posts/transport"
)

type App struct {
	server *server.Server
}

func New() (*App, error) {
	// Создание сервера
	srv := server.NewServer()

	// Создание подключения к базе данных
	db, err := db.NewDB(db.Config{
		Host:     viper.GetString("db.host"),     // получение хоста из конфига
		Port:     viper.GetString("db.port"),     // получение порта из конфига
		Username: viper.GetString("db.username"), // получение имени пользователя из конфига
		DBName:   viper.GetString("db.dbname"),   // получение имени базы данных из конфига
		SSLMode:  viper.GetString("db.sslmode"),  // получение режима SSL из конфига

		Password: os.Getenv("DB_PASSWORD"), // получение пароля из переменных окружения
	})

	// Проверка подключения
	if err != nil {
		logrus.Fatalf("ошибка при подключении к базе данных: [%s]\n", err)
	}

	// Инициализация зависимостей auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// Инициализация зависимостей posts
	postsRepo := postsR.NewRepository(db)
	postsService := postsS.NewService(postsRepo)
	postsHandler := postsT.NewHandler(postsService)

	// Регистрация маршрутов
	// Инициализация конфига роутера
	handlers.ConfigInit(srv.Router())

	authHandler.InitRouters(srv.Router())  // инициализация маршрутов для auth
	postsHandler.InitRouters(srv.Router()) // инициализация маршрутов для posts

	return &App{
		server: srv,
	}, nil
}

func (a *App) Run() error {
	return a.server.Run()
}
