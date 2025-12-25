package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vadimkiryanov/GO-CRUD/schema"
)

type ExamplePostgres struct {
	db *sqlx.DB
}

// NewExamplePostgres создает новый экземпляр структуры AuthPostgres
// db - это соединение с базой данных PostgreSQL
func NewExamplePostgres(db *sqlx.DB) *ExamplePostgres {
	// Возвращаем новый экземпляр AuthPostgres с установленным соединением к базе данных
	return &ExamplePostgres{db: db}
}

// Создает нового элемента в базе данных
func (repository *ExamplePostgres) CreateElement(element schema.ExampleSchema) (int, error) {
	// Переменная для хранения ID нового пользователя
	var id int

	// Формируем SQL запрос для вставки данных
	// $1, $2, $3 - это параметры, которые будут безопасно подставлены
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", examplesTable)

	// Выполняем запрос с данными пользователя
	// QueryRow используется, так как мы ожидаем только одну строку в ответе
	row := repository.db.QueryRow(query, element.Name)

	// Пытаемся получить ID созданного пользователя
	// Если произошла ошибка (например, дубликат username), возвращаем её
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	// Возвращаем ID нового пользователя и nil как ошибку
	return id, nil
}

// Получает элемент из базы данных по его name
func (repos *ExamplePostgres) GetElement(name string) (schema.ExampleSchema, error) {
	var exampleFromDb schema.ExampleSchema

	// Формируем SQL запрос
	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1", examplesTable)
	err := repos.db.Get(&exampleFromDb, query, name)

	return exampleFromDb, err // Возвращаем полученного пользователя и возможную ошибку
}
