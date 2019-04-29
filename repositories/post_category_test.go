package repositories

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestShouldFindPostCategoryById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name", "color"}).AddRow(1, "AWS", "#111111")
	mock.ExpectQuery("^SELECT (.+) FROM post_categories WHERE id = .*").WillReturnRows(row)

	postCategory, err := FindPostCategoryById(db, 1)
	if err != nil {
		t.Fatalf("Cannot get post_category: %s", err)
	}

	data := PostCategory{
		Id:    1,
		Name:  "AWS",
		Color: "#111111",
	}
	if postCategory != data {
		t.Fatalf("Fail expected: %v, got: %v", data, postCategory)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
