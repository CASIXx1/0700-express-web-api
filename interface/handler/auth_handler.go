package handler

import (
	"0700-express-web-api/interface/repository"

	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Handler struct {
	authRepository *repository.AuthRepository
	jwtSecret      string
}

func CreateHandler(authRepository *repository.AuthRepository, jwtSecret string) *Handler {
	return &Handler{
		authRepository: authRepository,
		jwtSecret:      jwtSecret,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
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

	result, err := handler.login(login.Email, login.Password)
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

func (handler *Handler) login(email, password string) (*loginResponse, error) {
	user, err := handler.authRepository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := handler.generateToken(user.ID.String(), "access", time.Now().Add(time.Hour))
	if err != nil {
		return nil, err
	}

	refreshToken, err := handler.generateToken(user.ID.String(), "refresh", time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, err
	}

	return &loginResponse{
		UUID:         user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (handler *Handler) generateToken(userId string, tokenType string, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userId,
		"type": tokenType,
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(handler.jwtSecret))
}
