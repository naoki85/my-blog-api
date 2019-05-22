package repositories

import (
	"database/sql"
	"time"
)

func GetPostsCount(db *sql.DB) (count int, err error) {
	nowTime := time.Now()
	query := "SELECT COUNT(*) FROM posts WHERE active = 1 AND published_at <= ? LIMIT 1"
	err = db.QueryRow(query, nowTime).Scan(&count)
	return
}
