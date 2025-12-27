package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vadimkiryanov/GO-CRUD/schema"
)

// signUp обрабатывает регистрацию пользователя путем привязки JSON-данных к структуре User.
// Возвращает ответ с ошибкой bad request, если привязка JSON не удалась.
func (handler *Handler) signUp(ctx *gin.Context) {
	// Когда пользователь отправляет запрос на регистрацию:
	var input schema.User // Создается пустая структура для данных пользователя

	// Пытаемся прочитать JSON из запроса и записать в структуру input todo.User
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Пытаемся создать пользователя через сервисный слой
	_, err := handler.services.Authorization.CreateUser(input)
	if err != nil {
		fmt.Println(err)
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := handler.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		fmt.Println(err)
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Отправляем успешный ответ (200) с ID созданного пользователя
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"` // binding:"required" - это означает что...
	Password string `json:"password" binding:"required"` // ...это поле является обязательным

}

func (handler *Handler) signIn(ctx *gin.Context) {
	var input signInInput

	// Пытаемся прочитать JSON из запроса и записать в структуру input
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Пытаемся создать токен через сервисный слой
	token, err := handler.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Отправляем успешный ответ (200)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
