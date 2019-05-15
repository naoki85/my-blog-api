package usecase

import (
	"github.com/naoki85/my_blog_api/models"
	"testing"
)

type MockRecommendedBookRepository struct {
}

func (repo *MockRecommendedBookRepository) All(int) (models.RecommendedBooks, error) {
	recommendedBooks := models.RecommendedBooks{
		models.RecommendedBook{Id: 1},
		models.RecommendedBook{Id: 2},
		models.RecommendedBook{Id: 3},
		models.RecommendedBook{Id: 4},
	}
	return recommendedBooks, nil
}

func TestShouldFindAllRecommendedBooks(t *testing.T) {
	repo := new(MockRecommendedBookRepository)
	interactor := RecommendedBookInteractor{
		RecommendedBookRepository: repo,
	}
	recommendedBooks, err := interactor.RecommendedBookRepository.All(4)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if len(recommendedBooks) != 4 {
		t.Fatalf("Fail expected: 4, got: %v", len(recommendedBooks))
	}
}
