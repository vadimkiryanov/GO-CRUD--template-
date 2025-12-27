package service

import (
	"github.com/vadimkiryanov/GO-CRUD/internal/repository"
	"github.com/vadimkiryanov/GO-CRUD/schema"
)

// AuthService структура для работы с аутентификацией
type ExampleService struct {
	repository repository.Example // Интерфейс для работы с хранилищем данных
}

// NewExampleService создает новый экземпляр сервиса аутентификации
func NewExampleService(repository repository.Example) *ExampleService {
	return &ExampleService{repository: repository}
}

// CreateElement создает нового пользователя
func (service *ExampleService) CreateElement(element schema.ExampleSchema) (int, error) {
	// Делегируем создание пользователя репозиторию
	return service.repository.CreateElement(element)
}

func (service *ExampleService) GetElement(name string) (schema.ExampleSchema, error) {
	return service.repository.GetElement(name)
}
