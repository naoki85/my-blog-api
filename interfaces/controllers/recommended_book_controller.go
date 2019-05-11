package controllers

import (
	"github.com/naoki85/my_blog_api/interfaces/database"
	"github.com/naoki85/my_blog_api/usecase"
	"strconv"
)

type RecommendedBookController struct {
	Interactor usecase.RecommendedBookInteractor
}

func NewRecommendedBookController(sqlHandler database.SqlHandler) *RecommendedBookController {
	return &RecommendedBookController{
		Interactor: usecase.RecommendedBookInteractor{
			RecommendedBookRepository: &database.RecommendedBookRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *RecommendedBookController) Index(c Context) {
	limit, _ := strconv.Atoi(c.Param("limit"))
	recommendedBooks, err := controller.Interactor.RecommendedBookRepository.All(limit)
	if err != nil {
		c.JSON(404, err.Error())
		return
	}
	c.JSON(200, recommendedBooks)
}
