package repositories

import (
	"database/sql"
	"time"
)

type Post struct {
	Id             int
	PostCategoryId int
	Title          string
	Content        string
	PublishedAt    string
	PostCategory   PostCategory
}

func FindPostById(db *sql.DB, id int) (post Post, err error) {
	nowTime := time.Now()
	query := "SELECT id, post_category_id, title, content, published_at FROM posts WHERE id = ? AND active = 1 AND published_at <= ? LIMIT 1"
	err = db.QueryRow(query, id, nowTime).
		Scan(&post.Id, &post.PostCategoryId, &post.Title, &post.Content, &post.PublishedAt)
	return
}
