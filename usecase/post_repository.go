package usecase

import "github.com/naoki85/my_blog_api/models"

type PostRepository interface {
	All(int) (models.Posts, error)
	FindById(int) (models.Post, error)
}
