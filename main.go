package main

import (
	"./config"
	. "./repositories"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type api struct {
	db *sql.DB
}

type post struct {
	Id             int
	PostCategoryId int
	Title          string
	ImageUrl       string
	PublishedAt    string
	PostCategory   PostCategory
}

func (a *api) posts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var page int
	queryPageParam := r.URL.Query().Get("page")
	if queryPageParam == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(queryPageParam)
	}
	offset := 10 * (page - 1)

	nowTime := time.Now()
	query := "SELECT id, post_category_id, title, image_file_name, published_at FROM posts WHERE active = 1 AND published_at <= ?"
	query = query + " ORDER BY published_at DESC LIMIT 10 OFFSET ?"
	rows, err := a.db.Query(query, nowTime, offset)
	if err != nil {
		a.fail(w, "failed to fetch posts: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var posts []*post
	for rows.Next() {
		p := &post{}
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
		p.PublishedAt = strings.Split(p.PublishedAt, " ")[0]

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
		Posts     []*post
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
	post.PublishedAt = strings.Split(post.PublishedAt, " ")[0]
	a.ok(w, post)
}

func (a *api) recommendedBooks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	recommendedBooks, err := FindAllRecommendedBooks(a.db, ParamsForFindAll{Limit: 4})
	if err != nil {
		a.fail(w, "Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		RecommendedBooks []*RecommendedBook
	}{recommendedBooks}

	a.ok(w, data)
}

func (a *api) handleOption(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if len(os.Args) > 1 {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3035")
	} else {
		w.Header().Add("Access-Control-Allow-Origin", "http://naoki85.me")
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := httprouter.New()
	app := &api{db: db}

	router.OPTIONS("/*path", app.handleOption)
	router.GET("/posts/:id", app.postById)
	router.GET("/posts", app.posts)
	router.GET("/recommended_books", app.recommendedBooks)
	http.ListenAndServe(":8080", router)
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
