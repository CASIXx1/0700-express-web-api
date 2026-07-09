package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Handler struct {
	db        *gorm.DB
	jwtSecret string
}

func CreateHandler(db *gorm.DB, jwtSecret string) *Handler {
	return &Handler{
		db:        db,
		jwtSecret: jwtSecret,
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

type user struct {
	ID       string `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	Status   string `gorm:"column:status"`
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
	var user user

	if err := handler.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, ErrInvalidCredentials
	}

	return &loginResponse{
		UUID:         user.ID,
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
	}, nil
}
