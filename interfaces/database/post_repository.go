package database

import (
	"github.com/naoki85/my_blog_api/models"
	"fmt"
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
		fmt.Println(error.Error())
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
