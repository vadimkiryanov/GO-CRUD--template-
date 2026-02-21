package transport

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vadimkiryanov/GO-CRUD/internal/core"
	"github.com/vadimkiryanov/GO-CRUD/internal/features/comments/service"
)

type Handler struct {
	services *service.CommentsService
}

func NewHandler(services *service.CommentsService) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouters(router *gin.Engine) *gin.Engine {
	comments := router.Group("/comments")
	{
		comments.GET("/post/:postId", h.getCommentsByPostId) // получение всех комментариев поста
		comments.POST("/create", h.createComment)            // создание комментария
		comments.PUT("/update/:id", h.updateComment)         // обновление комментария
		comments.DELETE("/delete/:id", h.deleteComment)      // удаление комментария
	}

	return router
}

func (handler *Handler) getCommentsByPostId(ctx *gin.Context) {
	postId := ctx.Param("postId")

	if postId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "post id is required"})
		return
	}

	var id int
	_, err := fmt.Sscanf(postId, "%d", &id)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	comments, err := handler.services.GetCommentsByPostId(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"comments": comments})
}

func (handler *Handler) createComment(ctx *gin.Context) {
	var input CommentsDto
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

	id, err := handler.services.CreateComment(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func (handler *Handler) deleteComment(ctx *gin.Context) {
	userId, _, err := core.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	commentId := ctx.Param("id")

	if commentId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "comment id is required"})
		return
	}

	var id int
	_, err = fmt.Sscanf(commentId, "%d", &id)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	err = handler.services.DeleteComment(id, userId)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "comment not found or access denied"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}

func (handler *Handler) updateComment(ctx *gin.Context) {
	userId, _, err := core.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	commentId := ctx.Param("id")

	if commentId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "comment id is required"})
		return
	}

	var id int
	_, err = fmt.Sscanf(commentId, "%d", &id)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	var input CommentsDto
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserId = userId

	err = handler.services.UpdateComment(id, userId, input)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "comment not found or access denied"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "comment updated"})
}
