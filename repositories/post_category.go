package repositories

import (
	"database/sql"
)

type PostCategory struct {
	Id    int
	Name  string
	Color string
}

func FindPostCategoryById(db *sql.DB, id int) (post_category PostCategory, err error) {
	err = db.QueryRow("SELECT id, name, color FROM post_categories WHERE id = ?", id).
		Scan(&post_category.Id, &post_category.Name, &post_category.Color)
	return
}