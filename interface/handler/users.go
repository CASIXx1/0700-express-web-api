package handler

import (
	"0700-express-web-api/usecase"
	"log"
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
	userID, ok := userIDFromContext(request.Context())
	if !ok || userID == "" {
		WriteResponse(writer, http.StatusUnauthorized, ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	user, err := handler.userUsecase.Me(request.Context(), userID)
	if err != nil {
		log.Printf("failed to get user: %v", err)

		WriteResponse(writer, http.StatusUnauthorized, ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	WriteResponse(writer, http.StatusOK, normalResponse[meResponse]{
		Data: meResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Status:   user.Status,
		},
	})
}
