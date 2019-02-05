package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type api struct {
	db *sql.DB
}

type post struct {
	ID          int
	Title       string
	Content     string
	PublishedAt string
}

func (a *api) posts(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.Query("SELECT id, title, content, published_at FROM posts")
	if err != nil {
		a.fail(w, "failed to fetch posts: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var posts []*post
	for rows.Next() {
		p := &post{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.PublishedAt); err != nil {
			a.fail(w, "failed to scan post: "+err.Error(), 500)
			return
		}
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

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(0.0.0.0:3306)/book_recorder_development")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	app := &api{db: db}
	http.HandleFunc("/posts", app.posts)
	http.ListenAndServe(":8080", nil)
}

func (a *api) fail(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Error string
	}{Error: msg}

	resp, _ := json.Marshal(data)
	w.WriteHeader(status)
	w.Write(resp)
}

func (a *api) ok(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.fail(w, "oops something evil has happened", 500)
		return
	}
	w.Write(resp)
}
