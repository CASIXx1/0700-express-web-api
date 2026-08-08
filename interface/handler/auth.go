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
		writeResponse(writer, http.StatusBadRequest, errorResponse{
			Message: err.Error(),
		})
		return
	}

	result, err := handler.authUsecase.Login(ctx, loginRequest.Email, loginRequest.Password)
	if err != nil {
		writeResponse(writer, http.StatusUnauthorized, errorResponse{
			Message: err.Error(),
		})
		return
	}

	loginRes := loginResponse{
		UUID:         result.UUID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	writeResponse(writer, http.StatusOK, normalResponse[loginResponse]{
		Data: loginRes,
	})
}

func (handler *Handler) SignUp(writer http.ResponseWriter, Request *http.Request) {
	var signUp signUpRequest
	ctx := Request.Context()

	if err := json.NewDecoder(Request.Body).Decode(&signUp); err != nil {
		writeResponse(writer, http.StatusBadRequest, errorResponse{
			Message: err.Error(),
		})
		return
	}

	result, err := handler.authUsecase.SignUp(ctx, signUp.Username, signUp.Email, signUp.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			writeResponse(writer, http.StatusConflict, errorResponse{
				Message: err.Error(),
			})
			return
		}

		writeResponse(writer, http.StatusBadRequest, errorResponse{
			Message: err.Error(),
		})
		return
	}

	signUpRes := signUpResponse{
		UUID:         result.UUID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	writeResponse(writer, http.StatusOK, normalResponse[signUpResponse]{
		Data: signUpRes,
	})
}
