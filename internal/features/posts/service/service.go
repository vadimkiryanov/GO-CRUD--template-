package service

import (
	"github.com/vadimkiryanov/GO-CRUD/internal/core/domains"
	"github.com/vadimkiryanov/GO-CRUD/internal/features/posts/repository"
)

type PostsService struct {
	repository repository.PostsRepository
}

func NewService(repository repository.PostsRepository) *PostsService {
	return &PostsService{repository: repository}
}

func (service *PostsService) CreatePost(post domains.PostsDomain) (int, error) {
	return service.repository.CreatePost(post)
}

func (service *PostsService) GetPosts(userId int) ([]repository.PostModel, error) {
	return service.repository.GetPosts(userId)
}

func (service *PostsService) DeletePost(postId int, userId int) error {
	return service.repository.DeletePost(postId, userId)
}

func (service *PostsService) UpdatePost(postId int, userId int, post domains.PostsDomain) error {
	return service.repository.UpdatePost(postId, userId, post)
}
