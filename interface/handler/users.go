package handler

import (
	"0700-express-web-api/usecase"
	"net/http"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

type meResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (handler *UserHandler) Me(writer http.ResponseWriter, request *http.Request) {
	accessToken, err := bearerToken(request)
	if err != nil {
		writeResponse(writer, http.StatusUnauthorized, errorResponse{
			Message: "unauthorized",
		})
		return
	}

	user, err := handler.userUsecase.Me(request.Context(), accessToken)
	if err != nil {
		writeResponse(writer, http.StatusUnauthorized, errorResponse{
			Message: "unauthorized",
		})
		return
	}

	writeResponse(writer, http.StatusOK, normalResponse[meResponse]{
		Data: meResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Status:   user.Status,
		},
	})
}
