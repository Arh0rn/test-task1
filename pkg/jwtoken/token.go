package jwtoken

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

func GenerateToken(id int, email string, secret []byte, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ParseToken(tokenString string, secret []byte) (int, error) {
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return 0, jwt.ErrInvalidKey
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}
	return int(id), nil
}

func ExtractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", jwt.ErrInvalidKey
	}

	parts := strings.Fields(authHeader)
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1], nil
	}
	if len(parts) == 1 {
		return parts[0], nil
	}

	return "", errors.New("invalid token")
}
