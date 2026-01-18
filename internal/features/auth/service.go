package auth

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vadimkiryanov/GO-CRUD/internal/core"
)

const (
	// Соль для усиления безопасности хеша
	salt     = "1a2b3c4dasdasdasd78921928"
	tokenTTL = 12 * time.Hour
)

type Repository interface {
	CreateUser(user User) (int, error)
	GetUser(username, password string) (User, error)
}

// AuthService структура для работы с аутентификацией
type AuthService struct {
	repository Repository // Интерфейс для работы с хранилищем данных
}

// NewAuthService создает новый экземпляр сервиса аутентификации
func NewService(repository Repository) *AuthService {
	return &AuthService{repository: repository}
}

// CreateUser создает нового пользователя
func (service *AuthService) CreateUser(user User) (int, error) {
	// Генерируем хеш пароля перед сохранением
	user.Password = generatePasswordHash(user.Password)
	// Делегируем создание пользователя репозиторию
	return service.repository.CreateUser(user)
}

// GenerateToken создает JWT токен для аутентифицированного пользователя
func (service *AuthService) GenerateToken(username, password string) (string, error) {
	// Получаем пользователя из базы данных, предварительно хешируя пароль
	// для сравнения с хешем в БД
	user, err := service.repository.GetUser(username, generatePasswordHash(password))
	if err != nil {
		// Если пользователь не найден или пароль неверный, возвращаем ошибку
		return "", err
	}

	// Создаем новый JWT токен с использованием алгоритма HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &core.TokenClaims{
		// StandardClaims содержит стандартные поля JWT
		StandardClaims: jwt.StandardClaims{
			// Устанавливаем время истечения токена (tokenTTL определено где-то в константах)
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			// Устанавливаем время создания токена
			IssuedAt: time.Now().Unix(),
		},
		// Добавляем ID пользователя в claims токена
		UserId:   user.Id,
		Username: user.Username,
	})

	// Подписываем токен секретным ключом и возвращаем его в виде строки
	// signInKey - это секретный ключ, определенный где-то в константах
	return token.SignedString([]byte(core.SignInKey))
}

// generatePasswordHash создает хеш пароля с использованием SHA1 и соли
func generatePasswordHash(password string) string {

	// Создаем новый SHA1 хеш
	hash := sha1.New()
	// Записываем пароль в хеш
	hash.Write([]byte(password))

	// Возвращаем хеш в виде шестнадцатеричной строки,
	// добавляя соль к результату
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
