package database

import (
	"github.com/naoki85/my_blog_api/models"
)

type RecommendedBookRepository struct {
	SqlHandler
}

func (repo *RecommendedBookRepository) All(limit int) (recommendedBooks models.RecommendedBooks, err error) {
	query := "SELECT id, link, image_url, button_url FROM recommended_books"
	query = query + " ORDER BY id DESC LIMIT ?"
	rows, err := repo.SqlHandler.Query(query, limit)
	defer rows.Close()
	if err != nil {
		return recommendedBooks, err
	}

	for rows.Next() {
		r := models.RecommendedBook{}
		if err := rows.Scan(&r.Id, &r.Link, &r.ImageUrl, &r.ButtonUrl); err != nil {
			return recommendedBooks, err
		}

		recommendedBooks = append(recommendedBooks, r)
	}
	return
}
