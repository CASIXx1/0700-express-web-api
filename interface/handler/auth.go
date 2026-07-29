package handler

import (
	"0700-express-web-api/usecase"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	authUsecase *usecase.AuthUsecase
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signUpRequest struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	EmailConfirmation    string `json:"email_confirmation"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type normalResponse[T any] struct {
	Data T `json:"data"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type loginResponse struct {
	UUID         string `json:"uuid"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type signUpResponse struct {
	UUID         string `json:"uuid"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewHandler(authUsecase *usecase.AuthUsecase) *Handler {
	return &Handler{
		authUsecase: authUsecase,
	}
}

func (handler *Handler) Login(writer http.ResponseWriter, Request *http.Request) {
	var loginRequest loginRequest
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

	loginResponse := loginResponse{
		UUID:         result.UUID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	writeResponse(writer, http.StatusOK, loginResponse)
}

func (handler *Handler) SignUp(writer http.ResponseWriter, Request *http.Request) {
	var signUp signUpRequest
	ctx := Request.Context()

	if err := json.NewDecoder(Request.Body).Decode(&signUp); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := handler.authUsecase.SignUp(ctx, signUp.Username, signUp.Email, signUp.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			writeErrorResponse(writer, http.StatusConflict, err)
			return
		}

		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	signUpResponse := signUpResponse{
		UUID:         result.UUID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	writeResponse(writer, http.StatusOK, signUpResponse)
}

func writeResponse(writer http.ResponseWriter, statusCode int, response any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	// ここがanyになってるので、気にしても仕方ない
	// json形式に変換するか
	// エラーの時も一緒に返せると良い
	json.NewEncoder(writer).Encode(map[string]any{
		"data": response,
	})
}

func writeErrorResponse(writer http.ResponseWriter, statusCode int, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(map[string]any{
		"error": err.Error(),
	})
}
