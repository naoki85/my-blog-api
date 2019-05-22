package controllers

import (
	"github.com/naoki85/my_blog_api/interfaces/database"
	"github.com/naoki85/my_blog_api/models"
	"github.com/naoki85/my_blog_api/usecase"
	"strconv"
	"strings"
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
	for _, p := range posts {
		p.PublishedAt = strings.Split(p.PublishedAt, "T")[0]
	}

	data := struct {
		Posts models.Posts
	}{posts}
	ok(w, data)
}

func (controller *PostController) Show(w ResponseWriter, r Request, p Params) {
	postId, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		fail(w, "Invalid Parameter", 400)
		return
	}
	post, err := controller.Interactor.FindById(postId)
	if err != nil || post.Id == 0 {
		fail(w, "Not Found", 404)
		return
	}
	ok(w, post)
}
