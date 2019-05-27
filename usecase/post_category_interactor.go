package usecase

import (
	"github.com/naoki85/my_blog_api/models"
)

type PostCategoryInteractor struct {
	PostCategoryRepository PostCategoryRepository
}

func (interactor *PostCategoryInteractor) FindById(id int) (models.PostCategory, error) {
	postCategory, err := interactor.PostCategoryRepository.FindById(id)
	if err != nil {
		return postCategory, err
	}
	return postCategory, err
}
