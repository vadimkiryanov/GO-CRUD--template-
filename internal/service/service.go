package service

import (
	"github.com/vadimkiryanov/GO-CRUD/internal/repository"
	"github.com/vadimkiryanov/GO-CRUD/schema"
)

type Example interface {
	CreateElement(schema.ExampleSchema) (int, error)
	GetElement(name string) (schema.ExampleSchema, error)
}

type Authorization interface {
	CreateUser(user schema.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (idUser int, err error)
}

type Service struct {
	Example
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Example:       NewExampleService(repos),
		Authorization: NewAuthService(repos),
	}
}
