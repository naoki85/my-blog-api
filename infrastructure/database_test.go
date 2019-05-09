package infrastructure

import (
	"testing"
)

func TestGetSqlHandler(t *testing.T) {
	_, err := NewSqlHandler()

	if err != nil {
		t.Fatalf("Cannot connect to database: %s", err)
	}
}
