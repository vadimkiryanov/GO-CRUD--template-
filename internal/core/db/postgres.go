package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// таблицы экспортируется в рамках пакета repository
// названия таблиц такие же, как в файлах миграций
const (
	examplesTable = "example"
	usersTable    = "users"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// NewDB создает новое подключение к базе данных PostgreSQL используя предоставленную конфигурацию.
// Функция пытается открыть соединение с указанными параметрами и выполняет ping для проверки подключения.
//
// Параметры:
//   - cfg Config: структура конфигурации, содержащая параметры подключения (хост, порт, имя пользователя,
//     имя базы данных, пароль и режим SSL)
//
// Возвращает:
//   - *sqlx.DB: указатель на объект подключения к базе данных
//   - error: ошибку, если не удалось установить соединение или выполнить ping
func NewDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
