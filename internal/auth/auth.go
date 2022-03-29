package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

var ErrInvalidToken = errors.New("token is invalid")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(secret, username string, expiresIn time.Duration) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresIn).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "could not create signed string")
	}

	return tokenString, nil
}

func IsAuthorized(secret, tokenString string) (bool, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, errors.Wrap(err, "could not parse claims")
	}

	if !token.Valid {
		return false, ErrInvalidToken
	}

	return true, nil
}
