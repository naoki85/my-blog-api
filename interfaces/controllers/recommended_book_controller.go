package controllers

import (
	"github.com/naoki85/my_blog_api/interfaces/database"
	"github.com/naoki85/my_blog_api/models"
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

func (controller *RecommendedBookController) Index(w ResponseWriter, r Request, p Params) {
	limit, _ := strconv.Atoi(p.ByName("limit"))
	recommendedBooks, err := controller.Interactor.RecommendedBookRepository.All(limit)
	if err != nil {
		fail(w, err.Error(), 404)
		return
	}

	data := struct {
		RecommendedBooks models.RecommendedBooks
	}{recommendedBooks}
	ok(w, data)
}
