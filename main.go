package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/naoki85/my_blog_api/infrastructure"
	"github.com/naoki85/my_blog_api/interfaces/controllers"
	"net/http"
	"os"
)

func main() {
	router := httprouter.New()

	router.OPTIONS("/*path", HandleOption)
	router.GET("/all_posts", AllPosts)
	router.GET("/posts/:id", NewGetPost)
	router.GET("/posts", NewPostsIndex)
	router.GET("/recommended_books", NewRecommendedBooks)
	http.ListenAndServe(":8080", router)
}

func HandleOption(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
