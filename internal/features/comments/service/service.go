package service

import (
	"github.com/vadimkiryanov/GO-CRUD/internal/core/domains"
	"github.com/vadimkiryanov/GO-CRUD/internal/features/comments/repository"
)

type CommentsService struct {
	repository repository.CommentsRepository
}

func NewService(repository repository.CommentsRepository) *CommentsService {
	return &CommentsService{repository: repository}
}

func (service *CommentsService) CreateComment(comment domains.CommentsDomain) (int, error) {
	return service.repository.CreateComment(comment)
}

func (service *CommentsService) GetCommentsByPostId(postId int) ([]repository.CommentModel, error) {
	return service.repository.GetCommentsByPostId(postId)
}

func (service *CommentsService) DeleteComment(commentId int, userId int) error {
	return service.repository.DeleteComment(commentId, userId)
}

func (service *CommentsService) UpdateComment(commentId int, userId int, comment domains.CommentsDomain) error {
	return service.repository.UpdateComment(commentId, userId, comment)
}
