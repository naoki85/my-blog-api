package database

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type MockSqlHandler struct {
	Conn *sql.DB
	Mock sqlmock.Sqlmock
}

func NewMockSqlHandler() (*MockSqlHandler, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err.Error())
	}

	if err != nil {
		panic(err.Error)
	}
	sqlHandler := new(MockSqlHandler)
	sqlHandler.Conn = db
	sqlHandler.Mock = mock
	return sqlHandler, err
}

func (handler *MockSqlHandler) ResistMock(statement string, fields []string, args ...interface{}) {
	rows := sqlmock.NewRows(fields).
		AddRow(1, "http://naoki85.test", "http://naoki85.test", "http://naoki85.test/button").
		AddRow(2, "http://naoki85.test", "http://naoki85.test/image", "http://naoki85.test/button").
		AddRow(3, "http://naoki85.test", "http://naoki85.test/image", "http://naoki85.test/button").
		AddRow(4, "http://naoki85.test", "http://naoki85.test/image", "http://naoki85.test/button")
	handler.Mock.ExpectQuery(statement).WillReturnRows(rows)
}

func (handler *MockSqlHandler) ResistMockForPosts(statement string, fields []string, args ...interface{}) {
	rows := sqlmock.NewRows(fields).
		AddRow(1, 1, "test title 1", "test content 1", "image_1", "2019-01-01 00:00:00").
		AddRow(2, 1, "test title 2", "test content 2", "image_2", "2019-01-02 00:00:00").
		AddRow(3, 1, "test title 3", "test content 3", "image_3", "2019-01-03 00:00:00").
		AddRow(4, 1, "test title 4", "test content 4", "image_4", "2019-01-04 00:00:00").
		AddRow(5, 1, "test title 5", "test content 5", "image_5", "2019-01-05 00:00:00")
	handler.Mock.ExpectQuery(statement).WillReturnRows(rows)
}

func (handler *MockSqlHandler) ResistMockForPost(statement string, fields []string, args ...interface{}) {
	rows := sqlmock.NewRows(fields).
		AddRow(1, 1, "test title 1", "test content 1", "image_1", "2019-01-01 00:00:00")
	handler.Mock.ExpectQuery(statement).WillReturnRows(rows)
}

func (handler *MockSqlHandler) Execute(statement string, args ...interface{}) (Result, error) {
	res := SqlResult{}
	result, err := handler.Conn.Exec(statement, args...)
	if err != nil {
		return res, err
	}
	res.Result = result
	return res, err
}

func (handler *MockSqlHandler) Query(statement string, args ...interface{}) (Row, error) {
	rows, err := handler.Conn.Query(statement, args...)
	if err != nil {
		return new(SqlRow), err
	}
	row := new(SqlRow)
	row.Rows = rows
	return row, nil
}

type SqlResult struct {
	Result sql.Result
}

func (r SqlResult) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

func (r SqlResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

type SqlRow struct {
	Rows *sql.Rows
}

func (r SqlRow) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r SqlRow) Next() bool {
	return r.Rows.Next()
}

func (r SqlRow) Close() error {
	return nil
}
