package controllers

import (
	"github.com/naoki85/my_blog_api/interfaces/database"
	"github.com/naoki85/my_blog_api/models"
	"github.com/naoki85/my_blog_api/usecase"
	"net/http"
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
			PostCategoryRepository: &database.PostCategoryRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *PostController) All(w ResponseWriter, r Request, p Params) {
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

func (controller *PostController) Index(w ResponseWriter, r *http.Request, p Params) {
	var page int
	queryPageParam := r.URL.Query().Get("page")
	if queryPageParam == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(queryPageParam)
	}
	posts, err := controller.Interactor.Index(page)
	if err != nil {
		fail(w, err.Error(), 404)
		return
	}

	var retPosts models.Posts

	if len(posts) == 0 {
		retPosts = models.Posts{}
	}

	for _, p := range posts {
		if p.ImageUrl == "" {
			p.ImageUrl = "https://s3-ap-northeast-1.amazonaws.com/bookrecorder-image/commons/default_user_icon.png"
		} else {
			p.ImageUrl = "http://d29xhtkvbwm2ne.cloudfront.net/" + p.ImageUrl
		}
		p.PublishedAt = strings.Split(p.PublishedAt, "T")[0]

		retPosts = append(retPosts, p)
	}

	count, err := controller.Interactor.GetPostsCount()
	if err != nil {
		fail(w, "failed to scan posts count: "+err.Error(), 500)
		return
	}
	totalPage := count / 10
	mod := count % 10
	if mod != 0 {
		totalPage = totalPage + 1
	}

	data := struct {
		TotalPage int
		Posts     models.Posts
	}{totalPage, retPosts}

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
