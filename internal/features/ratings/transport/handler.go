package transport

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vadimkiryanov/GO-CRUD/internal/core"
	"github.com/vadimkiryanov/GO-CRUD/internal/core/domains"
	"github.com/vadimkiryanov/GO-CRUD/internal/features/ratings/service"
)

type Handler struct {
	services *service.RatingService
}

// Инициализация обработчиков
func NewHandler(services *service.RatingService) *Handler {
	return &Handler{services: services}
}

// Инициализация роутеров
func (h *Handler) InitRouters(router *gin.Engine) *gin.Engine {
	ratings := router.Group("/ratings")
	{
		ratings.POST("/vote/:post_id", h.vote)                   // проголосовать за пост
		ratings.GET("/stats/:post_id", h.stats)                  // получить статистику по посту
		ratings.GET("/stats/:post_id/user", h.statsWithUserVote) // получить статистику с информацией о голосе пользователя
		ratings.DELETE("/vote/:post_id", h.deleteVote)           // удалить проголос пользователя за пост
	}

	return router
}

func (handler *Handler) vote(ctx *gin.Context) {
	userId, _, err := core.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postIdInt, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil || postIdInt <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var action domains.VoteAction
	if err := ctx.BindJSON(&action); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = handler.services.SetVote(userId, postIdInt, action)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "voted"})
}

func (handler *Handler) stats(ctx *gin.Context) {
	postIdInt, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil || postIdInt <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	stats, err := handler.services.GetStats(postIdInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (handler *Handler) statsWithUserVote(ctx *gin.Context) {
	userId, _, err := core.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postIdInt, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil || postIdInt <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	stats, err := handler.services.GetStatsWithUserVote(userId, postIdInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (handler *Handler) deleteVote(ctx *gin.Context) {
	userId, _, err := core.ParseToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postIdInt, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil || postIdInt <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	err = handler.services.DeleteVote(userId, postIdInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "vote deleted"})
}
