package models

type Post struct {
	Id             int
	PostCategoryId int
	Title          string
	Content        string
	ImageUrl       string
	PublishedAt    string
	PostCategory   PostCategory
}

type Posts []Post