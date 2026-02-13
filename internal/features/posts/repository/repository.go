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
	GetMyPosts(userId int) ([]PostModel, error)
	GetAllPosts() ([]PostModel, error)
	GetAllPostsWithPagination(params PaginationParams) (PaginatedResult, error)
	DeletePost(postId int, userId int) error
	UpdatePost(postId int, userId int, post domains.PostsDomain) error
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
func (repository *PostsPostgres) GetMyPosts(userId int) ([]PostModel, error) {
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

// Получает все посты из базы данных (без фильтрации по пользователю)
func (repository *PostsPostgres) GetAllPosts() ([]PostModel, error) {
	// Переменная для хранения списка постов
	var posts []PostModel

	// Формируем SQL запрос для получения всех постов
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
        JOIN users u ON p.user_id = u.id`, postsTable)

	// Выполняем запрос для получения всех постов
	// Queryx используется, так как мы ожидаем несколько строк в ответе
	err := repository.db.Select(&posts, query)
	// Возвращаем список постов и nil как ошибку
	return posts, err
}

// Получает все посты из базы данных с пагинацией
func (repository *PostsPostgres) GetAllPostsWithPagination(params PaginationParams) (PaginatedResult, error) {
	var posts []PostModel

	// Формируем SQL запрос для получения постов с лимитом и смещением
	query := fmt.Sprintf(`
        SELECT
            p.id,
            p.user_id,
            u.username as author,
            p.title,
            p.description,
            p.created_at,
            p.updated_at
        FROM %s p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC
        LIMIT $1 OFFSET $2`, postsTable)

	err := repository.db.Select(&posts, query, params.Limit, params.Offset)
	if err != nil {
		return PaginatedResult{}, err
	}

	// Получаем общее количество постов
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s p JOIN users u ON p.user_id = u.id`, postsTable)
	var total int64
	err = repository.db.Get(&total, countQuery)
	if err != nil {
		return PaginatedResult{}, err
	}

	// Вычисляем количество страниц
	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return PaginatedResult{
		Data:       posts,
		Total:      total,
		Page:       (params.Offset / params.Limit) + 1,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (repository *PostsPostgres) DeletePost(postId int, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", postsTable)
	_, err := repository.db.Exec(query, postId, userId)
	return err
}

func (repository *PostsPostgres) UpdatePost(postId int, userId int, post domains.PostsDomain) error {
	query := fmt.Sprintf("UPDATE %s SET title = $1, description = $2, updated_at = $3 WHERE id = $4 AND user_id = $5", postsTable)
	_, err := repository.db.Exec(query, post.Title, post.Description, time.Now(), postId, userId)
	return err
}
