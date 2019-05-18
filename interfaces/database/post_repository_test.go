package database

import (
	"testing"
)

func TestShouldFindAllPosts(t *testing.T) {
	mockSqlHandler, _ := NewMockSqlHandler()
	mockSqlHandler.ResistMockForPosts("^SELECT (.+) FROM posts .*", []string{"id", "post_category_id", "title", "content", "image_file_name", "published_at"})
	repo := PostRepository{
		SqlHandler: mockSqlHandler,
	}
	posts, err := repo.All(0)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if len(posts) != 5 {
		t.Fatalf("Fail expected: 5, got: %v", len(posts))
	}
	if posts[0].Title != "test title 1" {
		t.Fatalf("Fail expected: test title 1, got: %v", posts[0].Title)
	}
}
