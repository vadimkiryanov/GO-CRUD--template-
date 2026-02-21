package auth

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// таблицы экспортируется в рамках пакета repository
// названия таблиц такие же, как в файлах миграций
const (
	usersTable = "users"
)

type AuthPostgres struct {
	db *sqlx.DB
}

// NewAuthPostgres создает новый экземпляр структуры AuthPostgres
// db - это соединение с базой данных PostgreSQL
func NewRepository(db *sqlx.DB) *AuthPostgres {
	// Возвращаем новый экземпляр AuthPostgres с установленным соединением к базе данных
	return &AuthPostgres{db: db}
}

// Создает нового пользователя в базе данных
func (repository *AuthPostgres) CreateUser(user User) (int, error) {
	// Переменная для хранения ID нового пользователя
	var id int

	// Формируем SQL запрос для вставки данных
	// $1, $2, $3 - это параметры, которые будут безопасно подставлены
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)

	// Выполняем запрос с данными пользователя
	// QueryRow используется, так как мы ожидаем только одну строку в ответе
	row := repository.db.QueryRow(query, user.Name, user.Username, user.Password)

	// Пытаемся получить ID созданного пользователя
	// Если произошла ошибка (например, дубликат username), возвращаем её
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	// Возвращаем ID нового пользователя и nil как ошибку
	return id, nil
}

// Получает пользователя из базы данных по его username и паролю
func (repos *AuthPostgres) GetUser(username, password string) (User, error) {
	var userFromDb User

	// Формируем SQL запрос
	query := fmt.Sprintf("SELECT id, username FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := repos.db.Get(&userFromDb, query, username, password)

	return userFromDb, err // Возвращаем полученного пользователя и возможную ошибку
}
