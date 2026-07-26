package handler

import (
	"0700-express-web-api/interface/request"
	"0700-express-web-api/interface/response"
	"0700-express-web-api/usecase"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	authUsecase *usecase.AuthUsecase
}

func CreateHandler(authUsecase *usecase.AuthUsecase) *Handler {
	return &Handler{
		authUsecase: authUsecase,
	}
}

func (handler *Handler) Login(writer http.ResponseWriter, Request *http.Request) {
	var loginRequest request.LoginRequest
	ctx := Request.Context()

	if err := json.NewDecoder(Request.Body).Decode(&loginRequest); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := handler.authUsecase.Login(ctx, loginRequest.Email, loginRequest.Password)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	loginResponse := response.LoginResponse{
		UUID:         result.UUID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	writeResponse(writer, http.StatusOK, loginResponse)
}

func (handler *Handler) SignUp(writer http.ResponseWriter, Request *http.Request) {
	var signUp request.SignUpRequest
	ctx := Request.Context()

	if err := json.NewDecoder(Request.Body).Decode(&signUp); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := handler.authUsecase.SignUp(ctx, signUp.Username, signUp.Email, signUp.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			writer.WriteHeader(http.StatusConflict)
			return
		}

		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	signUpResponse := response.SignUpResponse{
		UUID:         result.UUID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	writeResponse(writer, http.StatusOK, signUpResponse)
}

func writeResponse(writer http.ResponseWriter, statusCode int, response any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(map[string]any{
		"data": response,
	})
}
