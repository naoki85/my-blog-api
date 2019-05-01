package repositories

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestShouldFindAllRecommendedBooks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "link", "image_url", "button_url"}).
		AddRow(1, "http://naoki85.test", "http://naoki85.test", "http://naoki85.test/button").
		AddRow(2, "http://naoki85.test", "http://naoki85.test/image", "http://naoki85.test/button").
		AddRow(3, "http://naoki85.test", "http://naoki85.test/image", "http://naoki85.test/button").
		AddRow(4, "http://naoki85.test", "http://naoki85.test/image", "http://naoki85.test/button")
	mock.ExpectQuery("^SELECT (.+) FROM recommended_books .*").WillReturnRows(rows)

	recommendedBooks, err := FindAllRecommendedBooks(db, ParamsForFindAll{limit: 4})
	if err != nil {
		t.Fatalf("Cannot get recommended_books: %s", err)
	}
	if len(recommendedBooks) != 4 {
		t.Fatalf("Fail expected: 4, got: %v", len(recommendedBooks))
	}
	if recommendedBooks[0].Link != "http://naoki85.test" {
		t.Fatalf("Fail expected: http://naoki85.test, got: %v", recommendedBooks[0].Link)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
