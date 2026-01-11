package main

import (
	_ "github.com/lib/pq" // Библиотека для работы с postgres, driver
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"github.com/vadimkiryanov/GO-CRUD/internal/app"
)

func main() {
	// Установка уровня логирования
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Инициализация конфига
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: [%s]\n", err)
	}

	// Инициализация переменных окружения
	if err := gotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: [%s]\n", err)
	}

	application, err := app.New()
	if err != nil {
		logrus.Fatal(err)
	}

	if err := application.Run(); err != nil {
		logrus.Fatal(err)
	}

	// Сервер запущен
	logrus.Info("Сервер запущен на порту: ", viper.GetString("port"))

	// Запуск сервера
	// если для viper.GetString key == неверное значение, то запустятся дефолтные настройки
	if err := application.Run(); err != nil {
		logrus.Fatalf("ошибка при запуске сервера: %s", err.Error())
	}

}

// initConfig инициализация конфига
func initConfig() error {
	viper.AddConfigPath("configs") // Папка с конфигами configs/
	viper.SetConfigName("config")  // Имя файла с конфигами configs/config.yaml

	return viper.ReadInConfig() // Чтение конфига

}
