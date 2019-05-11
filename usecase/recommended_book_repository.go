package usecase

import "github.com/naoki85/my_blog_api/models"

type RecommendedBookRepository interface {
	All(int) (models.RecommendedBooks, error)
}
