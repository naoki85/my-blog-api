package models

type RecommendedBook struct {
	Id        int
	Link      string
	ImageUrl  string
	ButtonUrl string
}

type RecommendedBooks []RecommendedBook
