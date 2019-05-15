package usecase

import (
	"github.com/naoki85/my_blog_api/models"
)

type RecommendedBookInteractor struct {
	RecommendedBookRepository RecommendedBookRepository
}

func (interactor *RecommendedBookInteractor) all(limit int) (models.RecommendedBooks, error) {
	recommendedBooks, err := interactor.RecommendedBookRepository.All(limit)
	return recommendedBooks, err
}
