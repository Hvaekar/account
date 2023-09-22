package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

func GenerateJWT(expTime time.Duration, payload interface{}, secretKey string) (*string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = time.Now().Add(expTime).Unix()

	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &tokenStr, nil
}

func ValidateJWT(token string, secretKey string) (interface{}, interface{}, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok || !tkn.Valid {
		return nil, nil, fmt.Errorf("invalid token claim")
	}

	return claims["sub"], claims["exp"], nil
}
