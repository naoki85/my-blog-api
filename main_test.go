package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	. "github.com/naoki85/my_blog_api/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (a *api) assertJSON(actual []byte, data interface{}, t *testing.T) {
	expected, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	if bytes.Compare(expected, actual) != 0 {
		t.Errorf("the expected json: %s is different from actual %s", expected, actual)
	}
}

func TestShouldGetPosts(t *testing.T) {
	dsn := fmt.Sprintf("root:root@tcp(0.0.0.0:3306)/book_recorder_test?parseTime=true&loc=Local")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_, _ = db.Exec("DELETE FROM posts")
		_, _ = db.Exec("DELETE FROM post_categories")
		db.Close()
	}()

	statement := "insert into posts (id, user_id, post_category_id, title, content, image_file_name, active, published_at, created_at, updated_at)"
	statement = statement + " values (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())"
	_, err = db.Exec(statement, 1, 1, 1, "test title 1", "test content 1", "image_1", 1, "2019-01-01 00:00:00")
	_, err = db.Exec(statement, 2, 1, 1, "test title 2", "test content 2", "image_2", 1, "2019-01-02 00:00:00")
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	statement = "insert into post_categories (id, name, color, created_at, updated_at) values (?, ?, ?, NOW(), NOW())"
	_, err = db.Exec(statement, 1, "Not categorized", "")
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

	// create app with mocked db, request and response to test
	app := &api{db}
	router := httprouter.New()
	router.GET("/posts", app.posts)

	req, err := http.NewRequest("GET", "http://localhost:8080/posts", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %v", w.Body)
	}

	data := struct {
		TotalPage int
		Posts     []*Post
	}{
		TotalPage: 1,
		Posts: []*Post{
			{Id: 2, PostCategoryId: 1, Title: "test title 2", ImageUrl: "http://d29xhtkvbwm2ne.cloudfront.net/image_2", PublishedAt: "2019-01-02", PostCategory: PostCategory{Id: 1, Name: "Not categorized", Color: ""}},
			{Id: 1, PostCategoryId: 1, Title: "test title 1", ImageUrl: "http://d29xhtkvbwm2ne.cloudfront.net/image_1", PublishedAt: "2019-01-01", PostCategory: PostCategory{Id: 1, Name: "Not categorized", Color: ""}},
		}}
	app.assertJSON(w.Body.Bytes(), data, t)
}

func TestShouldRespondWithErrorOnFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// create app with mocked db, request and response to test
	app := &api{db}
	router := httprouter.New()
	router.GET("/posts", app.posts)

	req, err := http.NewRequest("GET", "http://localhost/posts", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	// before we actually execute our api function, we need to expect required DB actions
	mock.ExpectQuery("^SELECT (.+) FROM posts .*").WillReturnError(fmt.Errorf("some error"))

	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Fatalf("expected status code to be 500, but got: %d", w.Code)
	}

	data := struct {
		Error string
	}{"failed to fetch posts: some error"}
	app.assertJSON(w.Body.Bytes(), data, t)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetPost(t *testing.T) {
	t.Skip("temporary skip")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// create app with mocked db, request and response to test
	app := &api{db}
	router := httprouter.New()
	router.GET("/posts/1", app.postById)

	req, err := http.NewRequest("GET", "http://localhost/posts/1", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	// before we actually execute our api function, we need to expect required DB actions
	rows := sqlmock.NewRows([]string{"id", "post_category_id", "title", "content", "published_at"}).
		AddRow(1, 1, "test title 1", "test content 1", "2019-01-01 00:00:00")

	mock.ExpectQuery("^SELECT (.+) FROM posts").WillReturnRows(rows)

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	data := Post{Id: 1, PostCategoryId: 1, Title: "test title 1", ImageUrl: "test.jpg", PublishedAt: "2019-01-01 00:00:00"}
	app.assertJSON(w.Body.Bytes(), data, t)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldRespondWithErrorOnNoPost(t *testing.T) {
	t.Skip("temporary skip")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// create app with mocked db, request and response to test
	app := &api{db}
	router := httprouter.New()
	router.GET("/posts/1", app.postById)

	req, err := http.NewRequest("GET", "http://localhost/posts/1", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	// before we actually execute our api function, we need to expect required DB actions
	rows := sqlmock.NewRows([]string{"id", "title", "content", "published_at"})

	mock.ExpectQuery("^SELECT (.+) FROM posts").WillReturnRows(rows)

	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
