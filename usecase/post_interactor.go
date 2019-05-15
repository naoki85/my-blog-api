package usecase

import (
	"github.com/naoki85/my_blog_api/models"
)

type PostInteractor struct {
	PostRepository PostRepository
}

func (interactor *PostInteractor) All(limit int) (models.Posts, error) {
	posts, err := interactor.PostRepository.All(limit)
	return posts, err
}
