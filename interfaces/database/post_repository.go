package database

import (
	"github.com/naoki85/my_blog_api/models"
	"time"
)

type PostRepository struct {
	SqlHandler
}

func (repo *PostRepository) All(limit int) (posts models.Posts, err error) {
	query := "SELECT id, post_category_id, title, content, image_file_name, published_at FROM posts WHERE active = 1 AND published_at <= ? ORDER BY published_at DESC"
	nowTime := time.Now()
	var rows Row
	var error error
	if limit > 0 {
		query = query + " LIMIT 10"
		rows, error = repo.SqlHandler.Query(query, nowTime, limit)
	} else {
		rows, error = repo.SqlHandler.Query(query, nowTime)
	}
	defer rows.Close()

	if error != nil {
		return posts, error
	}

	for rows.Next() {
		p := models.Post{}
		err := rows.Scan(&p.Id, &p.PostCategoryId, &p.Title, &p.Content, &p.ImageUrl, &p.PublishedAt)
		if err != nil {
			return posts, err
		}

		posts = append(posts, p)
	}
	return
}

func (repo *PostRepository) FindById(id int) (post models.Post, err error) {
	nowTime := time.Now()
	query := "SELECT id, post_category_id, title, content, image_file_name, published_at FROM posts WHERE id = ? AND active = 1 AND published_at <= ? LIMIT 1"
	rows, err := repo.SqlHandler.Query(query, id, nowTime)
	if err != nil {
		return post, err
	}

	for rows.Next() {
		err := rows.Scan(&post.Id, &post.PostCategoryId, &post.Title, &post.Content, &post.ImageUrl, &post.PublishedAt)
		if err != nil {
			return post, err
		}
		break
	}
	return
}
