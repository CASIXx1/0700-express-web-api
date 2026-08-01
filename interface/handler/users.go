package handler

import (
	"net/http"
	"strings"
)

type UserHandler struct {
}

type meRequest struct {
}

type meResponse struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Me(writer http.ResponseWriter, request *http.Request) {
	authorization := request.Header.Get("Authorization")

	if !strings.HasPrefix(authorization, "Bearer ") {
		http.Error(writer, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	//accessToken := strings.TrimPrefix(authorization, "Bearer ")

	// UserUsecaseを通して、ログインユーザーの情報を取得する
}
