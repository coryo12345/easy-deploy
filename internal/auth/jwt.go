package auth

import (
	"crypto/hmac"
	"crypto/sha512"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtBuilder interface {
	NewToken() (string, error)
	VerifyToken(token string) bool
}

type jwtBuilder struct {
	jwtKey []byte
}

func NewJwtBuilder() JwtBuilder {
	key := randomKey()
	return &jwtBuilder{
		jwtKey: key,
	}
}

func (b jwtBuilder) NewToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user",
		"iss": "easy-deploy",
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(b.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (b jwtBuilder) VerifyToken(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return b.jwtKey, nil
	})

	if err != nil {
		return false
	}
	return token.Valid
}

func randomKey() []byte {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 32)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	mac := hmac.New(sha512.New, []byte(string(s)))
	key := mac.Sum(nil)

	return key
}
