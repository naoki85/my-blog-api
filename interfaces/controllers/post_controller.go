package controllers

import (
	"github.com/naoki85/my_blog_api/interfaces/database"
	"github.com/naoki85/my_blog_api/models"
	"github.com/naoki85/my_blog_api/usecase"
)

type PostController struct {
	Interactor usecase.PostInteractor
}

func NewPostController(sqlHandler database.SqlHandler) *PostController {
	return &PostController{
		Interactor: usecase.PostInteractor{
			PostRepository: &database.PostRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *PostController) Index(w ResponseWriter, r Request, p Params) {
	limit := 0
	posts, err := controller.Interactor.PostRepository.All(limit)
	if err != nil {
		fail(w, err.Error(), 404)
		return
	}

	data := struct {
		Posts models.Posts
	}{posts}
	ok(w, data)
}
