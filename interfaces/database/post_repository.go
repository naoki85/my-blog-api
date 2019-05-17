package database

import (
	"github.com/naoki85/my_blog_api/models"
)

type PostRepository struct {
	SqlHandler
}

func (repo *PostRepository) All(limit int) (posts models.Posts, err error) {
	query := "SELECT id, post_category_id, title, image_file_name, published_at FROM posts WHERE active = 1 AND published_at <= ? ORDER BY published_at DESC"
	var rows Row
	var error error
	if limit > 0 {
		query = query + " LIMIT 10"
		rows, error = repo.SqlHandler.Query(query, limit)
	} else {
		rows, error = repo.SqlHandler.Query(query)
	}
	defer rows.Close()

	if error != nil {
		return posts, err
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
