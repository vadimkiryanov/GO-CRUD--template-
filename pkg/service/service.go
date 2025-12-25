package service

import (
	"github.com/vadimkiryanov/GO-CRUD/pkg/repository"
	"github.com/vadimkiryanov/GO-CRUD/schema"
)

type Example interface {
	CreateElement(schema.ExampleSchema) (int, error)
	GetElement(name string) (schema.ExampleSchema, error)
}

type Service struct {
	Example
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Example: NewExampleService(repos),
	}
}
