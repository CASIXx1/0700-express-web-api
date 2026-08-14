package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenService struct {
	secret string
}

func NewTokenService(secret string) *tokenService {
	return &tokenService{
		secret: secret,
	}
}

func (tokenService *tokenService) GenerateAccessToken(userID string) (string, error) {
	return generateToken(userID, "access", time.Now().Add(time.Hour), tokenService.secret)
}

func (tokenService *tokenService) GenerateRefreshToken(userID string) (string, error) {
	return generateToken(userID, "refresh", time.Now().Add(time.Hour*24*7), tokenService.secret)
}

func (tokenService *tokenService) VerifyAccessToken(accessToken string) (string, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenService.secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return "", errors.New("invalid token type")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token subject")
	}

	return userID, nil
}

func generateToken(userID string, tokenType string, expiresAt time.Time, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"type": tokenType,
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
