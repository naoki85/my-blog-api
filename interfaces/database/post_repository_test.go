package database

import (
	"database/sql"
	"fmt"
	"github.com/naoki85/my_blog_api/infrastructure"
	"testing"
)

func SetUpPostTest(t *testing.T) func() {
	dsn := fmt.Sprintf("root:root@tcp(0.0.0.0:3306)/book_recorder_test?parseTime=true&loc=Local")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	query := "INSERT INTO posts (user_id, title, content) VALUES(?, ?, ?)"
	db.Execute(query, 1, "test 1", "test 1 content")
	db.Execute(query, 1, "test 2", "test 2 content")
	db.Execute(query, 1, "test 3", "test 3 content")
	db.Execute(query, 1, "test 4", "test 4 content")
	db.Execute(query, 1, "test 5", "test 5 content")

	return func() {
		SqlHandler.Execute("DELETE FROM posts")
	}
}

func TestShouldFindAllPosts(t *testing.T) {
	tearDown := SetUpPostTest(t)
	defer tearDown()

	SqlHandler, _ := infrastructure.NewSqlHandler()
	repo := PostRepository{
		SqlHandler: SqlHandler,
	}
	posts, err := repo.All(0)
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if len(posts) != 4 {
		t.Fatalf("Fail expected: 4, got: %v", len(posts))
	}
	if posts[0].Title != "test 1" {
		t.Fatalf("Fail expected: http://naoki85.test, got: %v", posts[0].Title)
	}
}
