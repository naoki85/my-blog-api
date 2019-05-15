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
