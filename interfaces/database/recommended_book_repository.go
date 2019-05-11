package database

import (
	model "../../models"
)

type RecommendedBookRepository struct {
	SqlHandler
}

type ParamsForFindAll struct {
	Limit int
}

func (repo *RecommendedBookRepository) all(params ParamsForFindAll) (recommendedBooks model.RecommendedBooks, err error) {
	query := "SELECT id, link, image_url, button_url FROM recommended_books"
	query = query + " ORDER BY id DESC LIMIT ?"
	rows, err := repo.SqlHandler.Query(query, params.Limit)
	defer rows.Close()
	if err != nil {
		return recommendedBooks, err
	}

	for rows.Next() {
		r := model.RecommendedBook{}
		if err := rows.Scan(&r.Id, &r.Link, &r.ImageUrl, &r.ButtonUrl); err != nil {
			return recommendedBooks, err
		}

		recommendedBooks = append(recommendedBooks, r)
	}
	return
}
