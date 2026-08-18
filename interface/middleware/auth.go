package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type TokenVerifier interface {
	VerifyAccessToken(accessToken string) (string, error)
}

type contextKey string

const (
	userIDKey contextKey = "userID"
)

func Auth(tokenVerifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			accessToken, err := bearerToken(request)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusUnauthorized)
				return
			}

			userID, err := tokenVerifier.VerifyAccessToken(accessToken)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(request.Context(), userIDKey, userID)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func bearerToken(request *http.Request) (string, error) {
	authorization := request.Header.Get("Authorization")
	if authorization == "" {
		return "", errors.New("missing authorization")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid authorization")
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	if token == "" {
		return "", errors.New("invalid authorization")
	}

	return token, nil
}
