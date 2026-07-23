package handler

import (
	"0700-express-web-api/usecase"

	"encoding/json"
	"errors"
	"net/http"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserAlreadyExists = errors.New("user already exists")

type Handler struct {
	authUsecase *usecase.AuthUsecase
}

func CreateHandler(authUsecase *usecase.AuthUsecase) *Handler {
	return &Handler{
		authUsecase: authUsecase,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	EmailConfirmation    string `json:"email_confirmation"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
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

func (handler *Handler) Login(writer http.ResponseWriter, Request *http.Request) {
	var login loginRequest
	if err := json.NewDecoder(Request.Body).Decode(&login); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := handler.authUsecase.Login(login.Email, login.Password)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]any{
		"data": result,
	})
}

func (handler *Handler) SignUp(writer http.ResponseWriter, Request *http.Request) {
	var signUp SignUpRequest
	if err := json.NewDecoder(Request.Body).Decode(&signUp); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := handler.authUsecase.SignUp(signUp.Username, signUp.Email, signUp.Password)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			writer.WriteHeader(http.StatusConflict)
			return
		}

		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]any{
		"data": result,
	})
}
