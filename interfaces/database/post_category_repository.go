package database

import (
	"github.com/naoki85/my_blog_api/models"
)

type PostCategoryRepository struct {
	SqlHandler
}

func (repo *PostCategoryRepository) FindById(id int) (postCategory models.PostCategory, err error) {
	query := "SELECT id, name, color FROM post_categories WHERE id = ?"
	rows, err := repo.SqlHandler.Query(query, id)
	if err != nil {
		return postCategory, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&postCategory.Id, &postCategory.Name, &postCategory.Color)
		if err != nil {
			return postCategory, err
		}
		break
	}
	return
}
