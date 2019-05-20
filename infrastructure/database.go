package infrastructure

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/naoki85/my_blog_api/interfaces/database"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type SqlHandler struct {
	Conn *sql.DB
}

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

var instance *Config
var once sync.Once

func NewSqlHandler() (database.SqlHandler, error) {
	if len(os.Args) > 1 {
		Init(os.Args[1])
	} else {
		Init("")
	}
	c := Get()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Database)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error)
	}
	sqlHandler := new(SqlHandler)
	sqlHandler.Conn = conn
	return sqlHandler, err
}

func (handler *SqlHandler) Execute(statement string, args ...interface{}) (database.Result, error) {
	res := SqlResult{}
	result, err := handler.Conn.Exec(statement, args...)
	if err != nil {
		return res, err
	}
	res.Result = result
	return res, err
}

func (handler *SqlHandler) Query(statement string, args ...interface{}) (database.Row, error) {
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
	return r.Rows.Close()
}

func Get() *Config {
	return instance
}

func Init(e string) {
	env := e
	if e == "" {
		env = "production"
	} else if os.Getenv("ENV") == "test" {
		env = "test"
	}

	once.Do(func() {
		p, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		var filePath string
		if env == "test" {
			filePath = "/Users/yoneyamanaoki/go/src/github.com/naoki85/my_blog_api/config/database.yml"
		} else {
			filePath = filepath.Join(p, "/config/database.yml")
		}

		var conf map[string]Config
		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(buf, &conf)
		if err != nil {
			panic(err)
		}

		instance = &Config{
			Username: conf[env].Username,
			Password: conf[env].Password,
			Host:     conf[env].Host,
			Port:     conf[env].Port,
			Database: conf[env].Database,
		}
	})
}
