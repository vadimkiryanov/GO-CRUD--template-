package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vadimkiryanov/GO-CRUD/schema"
)

// CreateElement создает нового пользователя
// curl -X POST http://localhost:8000/test/element -H "Content-Type: application/json" -d '{"name": "vadim", "email": "vadim@gmail.com", "password": "123456"}'
func (h *Handler) CreateElement(ctx *gin.Context) {
	// Получаем данные из запроса
	var element schema.ExampleSchema
	if err := ctx.ShouldBindJSON(&element); err != nil {
		newErrorResponse(ctx, 400, err.Error())
		return
	}

	// Делегируем создание пользователя репозиторию
	id, err := h.services.CreateElement(element)
	if err != nil {
		newErrorResponse(ctx, 500, err.Error())
		return
	}

	// Возвращаем ID созданного пользователя
	ctx.JSON(200, gin.H{
		"id": id,
	})
}

// GetElement получает элемент из базы данных по его name
// curl -X GET http://localhost:8000/test/element/vadim
func (h *Handler) GetElement(ctx *gin.Context) {
	// Получаем данные из запроса
	name := ctx.Param("name")

	// Делегируем получение пользователя репозиторию
	element, err := h.services.GetElement(name)
	if err != nil {
		newErrorResponse(ctx, 500, err.Error())
		return
	}

	// Возвращаем полученного пользователя
	ctx.JSON(200, gin.H{
		"id":   element.Id,
		"name": element.Name,
	})
}
