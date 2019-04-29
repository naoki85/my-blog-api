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
	"time"
)

type api struct {
	db *sql.DB
}

type post struct {
	ID             int
	PostCategoryId int
	Title          string
	Content        string
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
	query := "SELECT id, post_category_id, title, content, published_at FROM posts WHERE active = 1 AND published_at <= ?"
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
		if err := rows.Scan(&p.ID, &p.PostCategoryId, &p.Title, &p.Content, &p.PublishedAt); err != nil {
			a.fail(w, "failed to scan post: "+err.Error(), 500)
			return
		}

		postCategory, err := FindPostCategoryById(a.db, p.PostCategoryId)
		if err != nil {
			a.fail(w, "failed to scan post_categories: "+err.Error(), 500)
			return
		}
		p.PostCategory = postCategory

		posts = append(posts, p)
	}

	if rows.Err() != nil {
		a.fail(w, "failed to read all posts: "+rows.Err().Error(), 500)
		return
	}

	data := struct {
		Posts []*post
	}{posts}

	a.ok(w, data)
}

func (a *api) postById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	postId := p.ByName("id")
	nowTime := time.Now()
	query := "SELECT id, post_category_id, title, content, published_at FROM posts WHERE id = ? AND active = 1 AND published_at <= ? LIMIT 1"
	rows, err := a.db.Query(query, postId, nowTime)
	if err != nil {
		a.fail(w, "failed to fetch posts: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var retPost post
	for rows.Next() {
		if err := rows.Scan(&retPost.ID, &retPost.PostCategoryId, &retPost.Title, &retPost.Content, &retPost.PublishedAt); err != nil {
			a.fail(w, "failed to scan post: "+err.Error(), 500)
			return
		}
	}
	if retPost.ID == 0 && retPost.Title == "" {
		a.fail(w, "Not found", 404)
		return
	}
	a.ok(w, retPost)
}

func (a *api) handleOption(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
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
