package transport

import (
	"fmt"
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
		auth.GET("/all", h.getAllPosts)          // получение всех постов
		auth.POST("/create", h.createPost)       // создание поста
		auth.PUT("/update/:id", h.updatePost)    // обновление поста
		auth.DELETE("/delete/:id", h.deletePost) // удаление поста
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

func (handler *Handler) deletePost(ctx *gin.Context) {
	userId, _, err := core.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postId := ctx.Param("id")

	// Валидация параметра ID
	if postId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "post id is required"})
		return
	}

	var id int
	_, err = fmt.Sscanf(postId, "%d", &id)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	err = handler.services.DeletePost(id, userId)
	if err != nil {
		// Возвращаем 404 если пост не найден или 403 если у пользователя нет доступа к этому посту
		ctx.JSON(http.StatusForbidden, gin.H{"error": "post not found or access denied"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

func (handler *Handler) updatePost(ctx *gin.Context) {
	userId, _, err := core.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postId := ctx.Param("id")

	// Валидация параметра ID
	if postId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "post id is required"})
		return
	}

	var id int
	_, err = fmt.Sscanf(postId, "%d", &id)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var input PostsDto
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Убедимся, что пользователь обновляет только свои посты
	input.UserId = userId

	err = handler.services.UpdatePost(id, userId, input)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "post not found or access denied"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "post updated"})
}
