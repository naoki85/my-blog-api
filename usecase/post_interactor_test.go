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

func (repo *MockPostRepository) Index(int) (models.Posts, error) {
	posts := models.Posts{
		models.Post{Id: 1, PostCategoryId: 1},
		models.Post{Id: 2, PostCategoryId: 1},
		models.Post{Id: 3, PostCategoryId: 1},
		models.Post{Id: 4, PostCategoryId: 1},
	}
	return posts, nil
}

func (repo *MockPostRepository) FindById(int) (models.Post, error) {
	post := models.Post{
		Id: 1,
	}
	return post, nil
}

func (repo *MockPostRepository) GetPostsCount() (int, error) {
	return 10, nil
}

func TestShouldFindAllPosts(t *testing.T) {
	interactor := initInteractor()
	posts, err := interactor.PostRepository.All(4)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if len(posts) != 4 {
		t.Fatalf("Fail expected: 4, got: %v", len(posts))
	}
}

func TestShouldPostsIndex(t *testing.T) {
	interactor := initInteractor()
	posts, err := interactor.Index(1)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if len(posts) != 4 {
		t.Fatalf("Fail expected: 4, got: %v", len(posts))
	}
	if posts[0].PostCategory.Name != "AWS" {
		t.Fatalf("Fail expected: AWS, got: %v", posts[0].PostCategory.Name)
	}
}

func TestShouldFindPostById(t *testing.T) {
	interactor := initInteractor()
	post, err := interactor.FindById(1)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if post.Id != 1 {
		t.Fatalf("Fail expected id: 1, got: %v", post)
	}
}

func TestShouldGetPostsCount(t *testing.T) {
	interactor := initInteractor()
	count, err := interactor.GetPostsCount()
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if count != 10 {
		t.Fatalf("Fail expected count: 10, got: %v", count)
	}
}

func initInteractor() PostInteractor {
	return PostInteractor{
		PostRepository:         new(MockPostRepository),
		PostCategoryRepository: new(MockPostCategoryRepository),
	}
}
