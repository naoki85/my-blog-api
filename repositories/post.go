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
	ImageUrl       string
	PublishedAt    string
	PostCategory   PostCategory
}

func FindPostById(db *sql.DB, id int) (post Post, err error) {
	nowTime := time.Now()
	query := "SELECT id, post_category_id, title, content, image_file_name, published_at FROM posts WHERE id = ? AND active = 1 AND published_at <= ? LIMIT 1"
	err = db.QueryRow(query, id, nowTime).
		Scan(&post.Id, &post.PostCategoryId, &post.Title, &post.Content, &post.ImageUrl, &post.PublishedAt)
	return
}

func GetPostsCount(db *sql.DB) (count int, err error) {
	nowTime := time.Now()
	query := "SELECT COUNT(*) FROM posts WHERE active = 1 AND published_at <= ? LIMIT 1"
	err = db.QueryRow(query, nowTime).Scan(&count)
	return
}
