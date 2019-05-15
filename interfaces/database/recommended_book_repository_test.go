package database

import (
	"testing"
)

func TestShouldFindAllRecommendedBooks(t *testing.T) {
	mockSqlHandler, _ := NewMockSqlHandler()
	mockSqlHandler.ResistMock("^SELECT (.+) FROM recommended_books .*", []string{"id", "link", "image_url", "button_url"})
	repo := RecommendedBookRepository{
		SqlHandler: mockSqlHandler,
	}
	recommendedBooks, err := repo.All(4)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if len(recommendedBooks) != 4 {
		t.Fatalf("Fail expected: 4, got: %v", len(recommendedBooks))
	}
	if recommendedBooks[0].Link != "http://naoki85.test" {
		t.Fatalf("Fail expected: http://naoki85.test, got: %v", recommendedBooks[0].Link)
	}
}
