package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vadimkiryanov/GO-CRUD/internal/core"
	"github.com/vadimkiryanov/GO-CRUD/internal/features/posts/service"
)

type Handler struct {
	services *service.PostsService
}

// Инициализация обработчиков
func NewHandler(services *service.PostsService) *Handler {
	return &Handler{services: services}
}

// Инициализация роутеров
func (h *Handler) InitRouters(router *gin.Engine) *gin.Engine {
	auth := router.Group("/posts")
	{
		auth.GET("/all", h.getAllPosts)    // получение всех постов
		auth.POST("/create", h.createPost) // получение всех постов
	}

	return router
}

func (handler *Handler) getAllPosts(ctx *gin.Context) {

	// Получаем пользователя из JWT
	userId, _, err := core.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Получаем все посты пользователя
	posts, err := handler.services.GetPosts(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"posts": posts})

}

func (handler *Handler) createPost(ctx *gin.Context) {
	var input PostsDto
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, userName, err := core.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	input.UserId = userId
	input.Author = userName

	id, err := handler.services.CreatePost(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})

}
