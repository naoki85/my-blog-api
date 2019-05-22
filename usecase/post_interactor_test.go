package usecase

import (
	"github.com/naoki85/my_blog_api/models"
	"testing"
)

type MockPostRepository struct {
}

func (repo *MockPostRepository) All(int) (models.Posts, error) {
	posts := models.Posts{
		models.Post{Id: 1},
		models.Post{Id: 2},
		models.Post{Id: 3},
		models.Post{Id: 4},
	}
	return posts, nil
}

func (repo *MockPostRepository) FindById(int) (models.Post, error) {
	post := models.Post{
		Id: 1,
	}
	return post, nil
}

func TestShouldFindAllPosts(t *testing.T) {
	repo := new(MockPostRepository)
	interactor := PostInteractor{
		PostRepository: repo,
	}
	posts, err := interactor.PostRepository.All(4)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if len(posts) != 4 {
		t.Fatalf("Fail expected: 4, got: %v", len(posts))
	}
}

func TestShouldFindPostById(t *testing.T) {
	repo := new(MockPostRepository)
	interactor := PostInteractor{
		PostRepository: repo,
	}
	post, err := interactor.PostRepository.FindById(1)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if post.Id != 1 {
		t.Fatalf("Fail expected id: 1, got: %v", post)
	}
}
