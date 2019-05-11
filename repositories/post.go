package repositories

import (
	"database/sql"
	. "github.com/naoki85/my_blog_api/models"
	"strings"
	"time"
)

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

func FindAllPosts(db *sql.DB) (posts []*Post, err error) {
	nowTime := time.Now()
	query := "SELECT id, post_category_id, title, content, image_file_name, published_at FROM posts WHERE active = 1 AND published_at <= ?"
	rows, err := db.Query(query, nowTime)
	defer rows.Close()

	for rows.Next() {
		p := &Post{}
		err := rows.Scan(&p.Id, &p.PostCategoryId, &p.Title, &p.Content, &p.ImageUrl, &p.PublishedAt)
		if err != nil {
			return posts, err
		}

		p.ImageUrl = "http://d29xhtkvbwm2ne.cloudfront.net/" + p.ImageUrl
		p.PublishedAt = strings.Split(p.PublishedAt, " ")[0]

		posts = append(posts, p)
	}
	return
}
