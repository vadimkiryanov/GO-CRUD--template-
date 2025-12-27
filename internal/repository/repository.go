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
type Authorization interface {
	CreateUser(user schema.User) (int, error)
	GetUser(username, password string) (schema.User, error)
}

type Repository struct {
	Example
	Authorization
}

// NewRepository создает новый экземпляр структуры Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Example:       NewExamplePostgres(db),
		Authorization: NewAuthPostgres(db),
	}
}
