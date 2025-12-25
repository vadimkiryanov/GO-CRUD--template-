package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/vadimkiryanov/GO-CRUD/schema"
)

// Example интерфейс репозитория для работы с авторизацией
type Example interface {
	CreateElement(element schema.ExampleSchema) (int, error)
	GetElement(name string) (schema.ExampleSchema, error)
}

type TodoItem interface {
}
type Repository struct {
	Example
}

// NewRepository создает новый экземпляр структуры Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Example: NewExamplePostgres(db),
	}
}
