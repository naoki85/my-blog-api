package repositories

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestShouldFindPostById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "post_category_id", "title", "content", "image_file_name", "published_at"}).
		AddRow(1, 1, "test title 1", "test content 1", "test_image_1", "2019-01-01 00:00:00")
	mock.ExpectQuery("^SELECT (.+) FROM posts").WillReturnRows(rows)

	post, err := FindPostById(db, 1)
	if err != nil {
		t.Fatalf("Cannot get post: %s", err)
	}

	data := Post{
		Id:             1,
		PostCategoryId: 1,
		Title:          "test title 1",
		Content:        "test content 1",
		ImageUrl:       "test_image_1",
		PublishedAt:    "2019-01-01 00:00:00",
	}
	if post != data {
		t.Fatalf("Fail expected: %v, got: %v", data, post)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetPostsCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(68)
	mock.ExpectQuery("^SELECT (.+) FROM posts .*").WillReturnRows(rows)

	count, err := GetPostsCount(db)
	if err != nil {
		t.Fatalf("Cannot get post: %s", err)
	}

	expected := 68
	if count != expected {
		t.Fatalf("Fail expected: %v, got: %v", expected, count)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldFindAllPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "post_category_id", "title", "content", "image_file_name", "published_at"}).
		AddRow(1, 1, "test title 1", "test content 1", "test_image_1", "2019-01-01 00:00:00").
		AddRow(2, 2, "test title 2", "test content 2", "test_image_2", "2019-01-01 00:00:00").
		AddRow(3, 3, "test title 3", "test content 3", "test_image_3", "2019-01-01 00:00:00")
	mock.ExpectQuery("^SELECT (.+) FROM posts").WillReturnRows(rows)

	posts, err := FindAllPosts(db)
	if err != nil {
		t.Fatalf("Cannot get posts: %s", err)
	}

	if len(posts) != 3 {
		t.Fatalf("Fail expected count: 3, got: %d", len(posts))
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
