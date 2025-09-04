package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/saiharsha/money-manager/internal/data"
)

var (
	ErrInvalidAuthenticationToken = errors.New("invalid or missing authentication token")
)

func (app *application) CreateToken(u *data.User, ttl time.Duration) (string, error) {
	expiresAt := time.Now().Add(ttl).Unix()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   u.ID,
		"role":  u.Role,
		"email": u.Email,
		"exp":   expiresAt,
		"iat":   time.Now().Add(ttl).Unix(),
	})

	token, err := claims.SignedString([]byte(app.config.secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (app *application) VerifyToken(tokenString string) (*data.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.config.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidAuthenticationToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	user := &data.User{
		ID:    int64(claims["sub"].(float64)),
		Email: fmt.Sprint(claims["email"]),
		Role:  fmt.Sprint(claims["role"]),
	}

	return user, nil
}
