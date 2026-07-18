package controller

import (
	"0700-express-web-api/ent"
	entuser "0700-express-web-api/ent/user"

	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Handler struct {
	dbClient  *ent.Client
	jwtSecret string
}

func CreateHandler(dbClient *ent.Client, jwtSecret string) *Handler {
	return &Handler{
		dbClient:  dbClient,
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
	user, err := handler.dbClient.User.
		Query().
		Where(entuser.EmailEQ(email)).
		First(context.Background())

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
