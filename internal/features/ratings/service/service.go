package service

import (
	"github.com/vadimkiryanov/GO-CRUD/internal/core/domains"
	"github.com/vadimkiryanov/GO-CRUD/internal/features/ratings/repository"
)

type RatingService struct {
	repo repository.RatingRepository
}

// NewService создает новый экземпляр сервиса RatingService
func NewService(repo repository.RatingRepository) *RatingService {
	return &RatingService{repo: repo}
}

func (s *RatingService) SetVote(userID, postID int, action domains.VoteAction) error {
	if err := s.repo.SetVote(userID, postID, action); err != nil {
		return err
	}
	// Опционально: обнови кэш/счетчики в posts
	return nil
}

func (s *RatingService) GetUserVote(userID, postID int) (*domains.Vote, error) {
	return s.repo.GetUserVote(userID, postID)
}

func (s *RatingService) GetStats(postID int) (*domains.RatingStats, error) {
	return s.repo.GetStats(postID)
}

func (s *RatingService) GetStatsWithUserVote(userID, postID int) (*domains.RatingStatsWithUserVote, error) {
	return s.repo.GetStatsWithUserVote(userID, postID)
}

func (s *RatingService) DeleteVote(userID, postID int) error {
	return s.repo.DeleteVote(userID, postID)
}
