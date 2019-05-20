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
	. "github.com/naoki85/my_blog_api/models"
	. "github.com/naoki85/my_blog_api/repositories"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type api struct {
	db *sql.DB
}

func (a *api) posts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var page int
	queryPageParam := r.URL.Query().Get("page")
	if queryPageParam == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(queryPageParam)
	}
	nowTime := time.Now()

	var rows *sql.Rows
	var err error

	query := "SELECT id, post_category_id, title, image_file_name, published_at FROM posts WHERE active = 1 AND published_at <= ? ORDER BY published_at DESC"
	query = query + " LIMIT 10 OFFSET ?"
	offset := 10 * (page - 1)
	rows, err = a.db.Query(query, nowTime, offset)

	if err != nil {
		a.fail(w, "failed to fetch posts: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		p := &Post{}
		if err := rows.Scan(&p.Id, &p.PostCategoryId, &p.Title, &p.ImageUrl, &p.PublishedAt); err != nil {
			a.fail(w, "failed to scan post: "+err.Error(), 500)
			return
		}

		postCategory, err := FindPostCategoryById(a.db, p.PostCategoryId)
		if err != nil {
			a.fail(w, "failed to scan post_categories: "+err.Error(), 500)
			return
		}
		p.PostCategory = postCategory
		p.ImageUrl = "http://d29xhtkvbwm2ne.cloudfront.net/" + p.ImageUrl
		p.PublishedAt = strings.Split(p.PublishedAt, "T")[0]

		posts = append(posts, p)
	}

	count, err := GetPostsCount(a.db)
	if err != nil {
		a.fail(w, "failed to scan posts count: "+err.Error(), 500)
		return
	}
	totalPage := count / 10
	mod := count % 10
	if mod != 0 {
		totalPage = totalPage + 1
	}

	if rows.Err() != nil {
		a.fail(w, "failed to read all posts: "+rows.Err().Error(), 500)
		return
	}

	data := struct {
		TotalPage int
		Posts     []*Post
	}{totalPage, posts}

	a.ok(w, data)
}

func (a *api) postById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	postId, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		a.fail(w, "Not Found", 404)
		return
	}
	post, err := FindPostById(a.db, postId)
	if err != nil {
		a.fail(w, "failed to fetch posts: "+err.Error(), 500)
		return
	}
	if post.Id == 0 && post.Title == "" {
		a.fail(w, "Not found", 404)
		return
	}
	post.PublishedAt = strings.Split(post.PublishedAt, "T")[0]
	post.ImageUrl = "http://d29xhtkvbwm2ne.cloudfront.net/" + post.ImageUrl
	a.ok(w, post)
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
	router.GET("/posts/:id", app.postById)
	router.GET("/posts", app.posts)
	router.GET("/recommended_books", NewRecommendedBooks)
	http.ListenAndServe(":8080", router)
}

func AllPosts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sqlHandler, _ := infrastructure.NewSqlHandler()
	postController := controllers.NewPostController(sqlHandler)
	postController.Index(w, r, p)
}

func NewRecommendedBooks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sqlHandler, _ := infrastructure.NewSqlHandler()
	recommendedBookController := controllers.NewRecommendedBookController(sqlHandler)
	recommendedBookController.Index(w, r, p)
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
