package repositories

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

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
