package usecase

import (
	"github.com/naoki85/my_blog_api/models"
	"strings"
)

type PostInteractor struct {
	PostRepository         PostRepository
	PostCategoryRepository PostCategoryRepository
}

func (interactor *PostInteractor) All(limit int) (models.Posts, error) {
	posts, err := interactor.PostRepository.All(limit)
	return posts, err
}

func (interactor *PostInteractor) Index(page int) (models.Posts, error) {
	posts, err := interactor.PostRepository.Index(page)
	var retPosts models.Posts
	for _, post := range posts {
		postCategory, err := interactor.PostCategoryRepository.FindById(post.PostCategoryId)
		if err != nil {
			continue
		}
		post.PostCategory = postCategory
		retPosts = append(retPosts, post)
	}
	return retPosts, err
}

func (interactor *PostInteractor) FindById(id int) (models.Post, error) {
	post, err := interactor.PostRepository.FindById(id)
	if err != nil {
		return post, err
	}
	post.PublishedAt = strings.Split(post.PublishedAt, "T")[0]
	post.ImageUrl = "http://d29xhtkvbwm2ne.cloudfront.net/" + post.ImageUrl
	return post, err
}
