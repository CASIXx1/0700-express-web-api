package middleware

import (
	"0700-express-web-api/interface/handler"
	"errors"
	"net/http"
	"strings"
)

type TokenVerifier interface {
	VerifyAccessToken(accessToken string) (string, error)
}

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

			ctx := handler.WithUserID(request.Context(), userID)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
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
