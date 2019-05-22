package usecase

import "github.com/naoki85/my_blog_api/models"

type PostCategoryRepository interface {
	FindById(int) (models.PostCategory, error)
}
