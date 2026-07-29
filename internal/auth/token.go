package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenGenerator struct {
	secret string
}

func NewTokenGenerator(secret string) *TokenGenerator {
	return &TokenGenerator{
		secret: secret,
	}
}

func (tokenGenerator *TokenGenerator) GenerateAccessToken(userId string) (string, error) {
	return generateToken(userId, "access", time.Now().Add(time.Hour), tokenGenerator.secret)
}

func (tokenGenerator *TokenGenerator) GenerateRefreshToken(userId string) (string, error) {
	return generateToken(userId, "refresh", time.Now().Add(time.Hour*24*7), tokenGenerator.secret)
}

func generateToken(userId string, tokenType string, expiresAt time.Time, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userId,
		"type": tokenType,
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
