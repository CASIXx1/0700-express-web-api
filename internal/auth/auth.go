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
	Token string `json:"token"`
}

type user struct {
	USERNAME string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler *Handler) Login(writer http.ResponseWriter, Request *http.Request) {
	var login loginRequest
	if err := json.NewDecoder(Request.Body).Decode(&login); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	//result, err := handler.login(login.Email, login.Password)
	//if err != nil {
	//	json.WriteError(writer, http.StatusUnauthorized, err)
	//	return
	//}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"status": "ok"})
}
