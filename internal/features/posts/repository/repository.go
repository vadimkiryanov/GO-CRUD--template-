package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vadimkiryanov/GO-CRUD/internal/core/domains"
)

// таблицы экспортируется в рамках пакета repository
// названия таблиц такие же, как в файлах миграций
const (
	postsTable = "posts"
)

type PostsPostgres struct {
	db *sqlx.DB
}

type PostsRepository interface {
	CreatePost(post domains.PostsDomain) (int, error)
	GetPosts(userId int) ([]PostModel, error)
}

// NewPostsPostgres создает новый экземпляр структуры PostsPostgres
// db - это соединение с базой данных PostgreSQL
func NewRepository(db *sqlx.DB) *PostsPostgres {
	// Возвращаем новый экземпляр PostsPostgres с установленным соединением к базе данных
	return &PostsPostgres{db: db}
}

// Создает нового поста в базе данных
func (repository *PostsPostgres) CreatePost(post domains.PostsDomain) (int, error) {
	// Переменная для хранения ID нового поста
	var id int

	// Формируем SQL запрос для вставки данных
	// $1, $2, $3 - это параметры, которые будут безопасно подставлены
	query := fmt.Sprintf("INSERT INTO %s (user_id, title, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", postsTable)

	// Выполняем запрос с данными поста
	// QueryRow используется, так как мы ожидаем только одну строку в ответе
	postFinal := PostModel{
		UserId:      post.UserId,
		Title:       post.Title,
		Description: post.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	row := repository.db.QueryRow(query, postFinal.UserId, postFinal.Title, postFinal.Description, postFinal.CreatedAt, postFinal.UpdatedAt)

	// Пытаемся получить ID созданного поста
	// Если произошла ошибка (например, дубликат user_id), возвращаем её
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	// Возвращаем ID нового поста и nil как ошибку
	return id, nil
}

// Получает все посты пользователя из базы данных
func (repository *PostsPostgres) GetPosts(userId int) ([]PostModel, error) {
	// Переменная для хранения списка постов
	var posts []PostModel

	// Формируем SQL запрос для получения всех постов пользователя
	// JOIN с users для получения имени автора
	query := fmt.Sprintf(`
        SELECT 
            p.id, 
            p.user_id, 
            u.username as author,  -- <-- имя пользователя из users.username
            p.title, 
            p.description, 
            p.created_at, 
            p.updated_at 
        FROM %s p 
        JOIN users u ON p.user_id = u.id 
        WHERE p.user_id = $1`, postsTable)

	// Выполняем запрос с данными пользователя
	// Queryx используется, так как мы ожидаем несколько строк в ответе
	err := repository.db.Select(&posts, query, userId)
	// Возвращаем список постов и nil как ошибку
	return posts, err
}
