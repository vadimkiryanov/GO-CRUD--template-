package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vadimkiryanov/GO-CRUD/internal/repository"
	"github.com/vadimkiryanov/GO-CRUD/schema"
)

const (
	// Соль для усиления безопасности хеша
	salt      = "1a2b3c4dasdasdasd78921928"
	tokenTTL  = 12 * time.Hour
	signInKey = "huidasui#81j12ASUidhqi81X"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

// AuthService структура для работы с аутентификацией
type AuthService struct {
	repository repository.Authorization // Интерфейс для работы с хранилищем данных
}

// NewAuthService создает новый экземпляр сервиса аутентификации
func NewAuthService(repository repository.Authorization) *AuthService {
	return &AuthService{repository: repository}
}

// ParseToken принимает токен доступа и возвращает ID пользователя и ошибку
func (service *AuthService) ParseToken(accsessToken string) (idUser int, err error) {
	// Парсим JWT токен с помощью jwt.ParseWithClaims
	// Эта функция проверяет подпись и декодирует данные токена
	token, err := jwt.ParseWithClaims(
		accsessToken,   // Сам токен доступа
		&tokenClaims{}, // Структура, в которую будут декодированы данные токена
		// Функция для проверки метода подписи и получения ключа
		func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что метод подписи токена - HMAC
			// ok будет false, если использован другой метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			// Возвращаем ключ для проверки подписи
			return []byte(signInKey), nil
		})

	// Если возникла ошибка при парсинге токена
	if err != nil {
		return 0, err
	}

	// Пытаемся привести claims к нашему типу tokenClaims
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	// Возвращаем ID пользователя из токена
	return claims.UserId, nil
}

// CreateUser создает нового пользователя
func (service *AuthService) CreateUser(user schema.User) (int, error) {
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		// StandardClaims содержит стандартные поля JWT
		jwt.StandardClaims{
			// Устанавливаем время истечения токена (tokenTTL определено где-то в константах)
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			// Устанавливаем время создания токена
			IssuedAt: time.Now().Unix(),
		},
		// Добавляем ID пользователя в claims токена
		user.Id,
	})

	// Подписываем токен секретным ключом и возвращаем его в виде строки
	// signInKey - это секретный ключ, определенный где-то в константах
	return token.SignedString([]byte(signInKey))
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
