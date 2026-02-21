package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vadimkiryanov/GO-CRUD/internal/core/domains"
)

const commentsTable = "comments"

type CommentsPostgres struct {
	db *sqlx.DB
}

type CommentsRepository interface {
	CreateComment(comment domains.CommentsDomain) (int, error)
	GetCommentsByPostId(postId int) ([]CommentModel, error)
	DeleteComment(commentId int, userId int) error
	UpdateComment(commentId int, userId int, comment domains.CommentsDomain) error
}

func NewRepository(db *sqlx.DB) *CommentsPostgres {
	return &CommentsPostgres{db: db}
}

func (repository *CommentsPostgres) CreateComment(comment domains.CommentsDomain) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (post_id, user_id, author, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", commentsTable)

	commentFinal := CommentModel{
		PostId:    comment.PostId,
		UserId:    comment.UserId,
		Author:    comment.Author,
		Content:   comment.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	row := repository.db.QueryRow(query, commentFinal.PostId, commentFinal.UserId, commentFinal.Author, commentFinal.Content, commentFinal.CreatedAt, commentFinal.UpdatedAt)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repository *CommentsPostgres) GetCommentsByPostId(postId int) ([]CommentModel, error) {
	var comments []CommentModel

	query := fmt.Sprintf(`
        SELECT
            c.id,
            c.post_id,
            c.user_id,
            c.author,
            c.content,
            c.created_at,
            c.updated_at
        FROM %s c
        WHERE c.post_id = $1
        ORDER BY c.created_at ASC`, commentsTable)

	err := repository.db.Select(&comments, query, postId)
	return comments, err
}

func (repository *CommentsPostgres) DeleteComment(commentId int, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", commentsTable)
	_, err := repository.db.Exec(query, commentId, userId)
	return err
}

func (repository *CommentsPostgres) UpdateComment(commentId int, userId int, comment domains.CommentsDomain) error {
	query := fmt.Sprintf("UPDATE %s SET content = $1, updated_at = $2 WHERE id = $3 AND user_id = $4", commentsTable)
	_, err := repository.db.Exec(query, comment.Content, time.Now(), commentId, userId)
	return err
}
