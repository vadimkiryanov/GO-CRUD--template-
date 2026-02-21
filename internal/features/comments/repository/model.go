package repository

import "time"

type CommentModel struct {
	ID        int       `db:"id"         json:"id"`
	PostId    int       `db:"post_id"    json:"post_id"`
	UserId    int       `db:"user_id"    json:"user_id"`
	Author    string    `db:"author"     json:"author"`
	Content   string    `db:"content"    json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
