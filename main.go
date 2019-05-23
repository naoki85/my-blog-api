package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/naoki85/my_blog_api/config"
	"github.com/naoki85/my_blog_api/infrastructure"
	"github.com/naoki85/my_blog_api/interfaces/controllers"
	"net/http"
	"os"
)

type api struct {
	db *sql.DB
}

func (a *api) handleOption(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if len(os.Args) > 1 {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3035")
	} else {
		w.Header().Add("Access-Control-Allow-Origin", "https://blog.naoki85.me")
	}
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Headers", "Origin")
	w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With")
	w.Header().Add("Access-Control-Allow-Headers", "Accept")
	w.Header().Add("Access-Control-Allow-Headers", "Accept-Language")
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	initEnv()
	c := config.Get()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := httprouter.New()
	app := &api{db: db}

	router.OPTIONS("/*path", app.handleOption)
	router.GET("/all_posts", AllPosts)
	router.GET("/posts/:id", NewGetPost)
	router.GET("/posts", NewPostsIndex)
	router.GET("/recommended_books", NewRecommendedBooks)
	http.ListenAndServe(":8080", router)
}

func AllPosts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sqlHandler, _ := infrastructure.NewSqlHandler()
	postController := controllers.NewPostController(sqlHandler)
	postController.All(w, r, p)
}

func NewPostsIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sqlHandler, _ := infrastructure.NewSqlHandler()
	postController := controllers.NewPostController(sqlHandler)
	postController.Index(w, r, p)
}

func NewRecommendedBooks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sqlHandler, _ := infrastructure.NewSqlHandler()
	recommendedBookController := controllers.NewRecommendedBookController(sqlHandler)
	recommendedBookController.Index(w, r, p)
}

func NewGetPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sqlHandler, _ := infrastructure.NewSqlHandler()
	postController := controllers.NewPostController(sqlHandler)
	postController.Show(w, r, p)
}

func (a *api) fail(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data := struct {
		Error string
	}{Error: msg}

	resp, _ := json.Marshal(data)
	w.WriteHeader(status)
	w.Write(resp)
}

func (a *api) ok(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.fail(w, "oops something evil has happened", 500)
		return
	}
	w.Write(resp)
}

func initEnv() {
	if len(os.Args) > 1 {
		config.Init(os.Args[1])
	} else {
		config.Init("")
	}
}
