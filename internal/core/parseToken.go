package core

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

const SignInKey = "huidasui#81j12ASUidhqi81X"

type TokenClaims struct {
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

// ParseToken принимает токен доступа и возвращает ID пользователя и ошибку
func ParseToken(accsessToken string) (idUser int, username string, err error) {
	// Парсим JWT токен с помощью jwt.ParseWithClaims
	// Эта функция проверяет подпись и декодирует данные токена
	token, err := jwt.ParseWithClaims(
		accsessToken,   // Сам токен доступа
		&TokenClaims{}, // Структура, в которую будут декодированы данные токена
		// Функция для проверки метода подписи и получения ключа
		func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что метод подписи токена - HMAC
			// ok будет false, если использован другой метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			// Возвращаем ключ для проверки подписи
			return []byte(SignInKey), nil
		})

	// Если возникла ошибка при парсинге токена
	if err != nil {
		return 0, "", err
	}

	// Пытаемся привести claims к нашему типу tokenClaims
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *tokenClaims")
	}

	// Возвращаем ID пользователя из токена
	return claims.UserId, claims.Username, nil
}
