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

func (service *PostsService) GetMyPosts(userId int) ([]repository.PostModel, error) {
	return service.repository.GetMyPosts(userId)
}

func (service *PostsService) GetAllPosts(sortBy string, sortOrder string) ([]repository.PostModel, error) {
	return service.repository.GetAllPosts(sortBy, sortOrder)
}

func (service *PostsService) DeletePost(postId int, userId int) error {
	return service.repository.DeletePost(postId, userId)
}

func (service *PostsService) UpdatePost(postId int, userId int, post domains.PostsDomain) error {
	return service.repository.UpdatePost(postId, userId, post)
}

func (service *PostsService) GetAllPostsWithPagination(params repository.PaginationParams) (repository.PaginatedResult, error) {
	return service.repository.GetAllPostsWithPagination(params)
}
