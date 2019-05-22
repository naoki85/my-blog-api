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
	"strconv"
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
	db, tearDown := InitTestDb()
	defer tearDown()

	statement := "insert into posts (id, user_id, post_category_id, title, content, image_file_name, active, published_at, created_at, updated_at)"
	statement = statement + " values (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())"
	_, err := db.Exec(statement, 1, 1, 1, "test title 1", "test content 1", "image_1", 1, "2019-01-01 00:00:00")
	_, err = db.Exec(statement, 2, 1, 1, "test title 2", "test content 2", "image_2", 1, "2019-01-02 00:00:00")
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	statement = "insert into post_categories (id, name, color, created_at, updated_at) values (?, ?, ?, NOW(), NOW())"
	_, err = db.Exec(statement, 1, "Not categorized", "")
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

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
	db, tearDown := InitTestDb()
	defer tearDown()

	statement := "insert into posts (id, user_id, post_category_id, title, content, image_file_name, active, published_at, created_at, updated_at)"
	statement = statement + " values (1, 1, 1, ?, ?, ?, ?, ?, NOW(), NOW())"
	_, err := db.Exec(statement, "test title 1", "test content 1", "image_1", 1, "2019-01-01 00:00:00")
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

	app := &api{db}
	router := httprouter.New()
	router.GET("/posts/:id", NewGetPost)

	req, err := http.NewRequest("GET", "http://localhost:8080/posts/1", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %v", w.Body)
	}

	data := Post{Id: 1, PostCategoryId: 1, Title: "test title 1", Content: "test content 1", ImageUrl: "http://d29xhtkvbwm2ne.cloudfront.net/image_1", PublishedAt: "2019-01-01"}
	app.assertJSON(w.Body.Bytes(), data, t)

	req, err = http.NewRequest("GET", "http://localhost:8080/posts/2", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Fatalf("expected status code to be 404, but got: %v", w.Body)
	}
}

func TestShouldGetAllPosts(t *testing.T) {
	db, tearDown := InitTestDb()
	defer tearDown()

	for i := 1; i <= 20; i++ {
		iToStr := strconv.Itoa(i)
		statement := "insert into posts (id, user_id, post_category_id, title, content, image_file_name, active, published_at, created_at, updated_at)"
		statement = statement + " values (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())"
		_, err := db.Exec(statement, i, 1, 1, "test title "+iToStr, "test content "+iToStr, "image_"+iToStr, 1, "2019-01-01 00:00:00")
		if err != nil {
			t.Fatalf("an error '%s' was not expected while creating request", err)
		}
	}

	statement := "insert into post_categories (id, name, color, created_at, updated_at) values (?, ?, ?, NOW(), NOW())"
	_, err := db.Exec(statement, 1, "Not categorized", "")
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

	router := httprouter.New()
	router.GET("/all_posts", AllPosts)

	req, err := http.NewRequest("GET", "http://localhost:8080/all_posts", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %v", w.Body)
	}

	var TestResponse struct {
		TotalPage int
		Posts     Posts
	}

	Bytes := []byte(w.Body.Bytes())
	json.Unmarshal(Bytes, &TestResponse)

	if len(TestResponse.Posts) != 20 {
		t.Fatalf("Fail expected: 20, got: %v", len(TestResponse.Posts))
	}
}

func TestShouldRespondWithErrorOnNoPost(t *testing.T) {
	db, tearDown := InitTestDb()
	defer tearDown()

	statement := "insert into posts (id, user_id, post_category_id, title, content, image_file_name, active, published_at, created_at, updated_at)"
	statement = statement + " values (1, 1, 1, ?, ?, ?, ?, ?, NOW(), NOW())"
	_, err := db.Exec(statement, "test title 1", "test content 1", "image_1", 1, "2019-01-01 00:00:00")
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

	router := httprouter.New()
	router.GET("/posts/:id", NewGetPost)

	req, err := http.NewRequest("GET", "http://localhost:8080/posts/2", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Fatalf("expected status code to be 404, but got: %v", w.Body)
	}
}

func InitTestDb() (*sql.DB, func()) {
	dsn := fmt.Sprintf("root:root@tcp(0.0.0.0:3306)/book_recorder_test?parseTime=true&loc=Local")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db, func() {
		_, _ = db.Exec("DELETE FROM posts")
		_, _ = db.Exec("DELETE FROM post_categories")
		db.Close()
	}
}
