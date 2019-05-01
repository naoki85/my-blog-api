package repositories

import (
	"database/sql"
)

type RecommendedBook struct {
	Id        int
	Link      string
	ImageUrl  string
	ButtonUrl string
}

type ParamsForFindAll struct {
	Limit int
}

func FindAllRecommendedBooks(db *sql.DB, params ParamsForFindAll) (recommendedBooks []*RecommendedBook, err error) {
	query := "SELECT id, link, image_url, button_url FROM recommended_books"
	query = query + " ORDER BY id DESC LIMIT ?"
	rows, err := db.Query(query, params.Limit)
	defer rows.Close()
	if err != nil {
		return recommendedBooks, err
	}

	for rows.Next() {
		r := &RecommendedBook{}
		if err := rows.Scan(&r.Id, &r.Link, &r.ImageUrl, &r.ButtonUrl); err != nil {
			return recommendedBooks, err
		}

		recommendedBooks = append(recommendedBooks, r)
	}
	return
}
