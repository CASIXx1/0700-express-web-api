package handler

import (
	"net/http"
	"strings"
)

func bearerToken(request *http.Request) (string, bool) {
	authorization := request.Header.Get("Authorization")
	if authorization == "" {
		return "", false
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", false
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	if token == "" {
		return "", false
	}

	return token, true
}
