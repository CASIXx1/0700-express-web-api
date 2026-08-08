package handler

import (
	"errors"
	"net/http"
	"strings"
)

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
