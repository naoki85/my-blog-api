package infrastructure

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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

func NewSqlHandler() (sqlHandler *SqlHandler, err error) {
	Init()
	c := Get()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Database)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return sqlHandler, err
	}
	defer conn.Close()
	sqlHandler = &SqlHandler{Conn: conn}
	fmt.Println(sqlHandler)
	return
}

func Get() *Config {
	return instance
}

func Init() {
	var env string
	env = "test"

	once.Do(func() {
		p, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		var filePath string
		filePath = p + "/../config/database.yml"

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
