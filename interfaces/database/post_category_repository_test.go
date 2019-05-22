package database

import (
	"testing"
)

func TestShouldFindPostCategory(t *testing.T) {
	mockSqlHandler, _ := NewMockSqlHandler()
	mockSqlHandler.ResistMockForPostCategory("^SELECT (.+) FROM post_categories (.+)", []string{"id", "name", "color"})
	repo := PostCategoryRepository{
		SqlHandler: mockSqlHandler,
	}
	postCategory, err := repo.FindById(1)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if postCategory.Name != "AWS" {
		t.Fatalf("Fail expected: AWS, got: %v", postCategory.Name)
	}
}
