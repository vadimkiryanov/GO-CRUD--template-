package repository

import "time"

type PostModel struct {
	ID          int       `db:"id"          json:"id"`
	Author      string    `db:"author"      json:"author"`
	UserId      int       `db:"user_id"     json:"user_id"`
	Title       string    `db:"title"       json:"title"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
}
