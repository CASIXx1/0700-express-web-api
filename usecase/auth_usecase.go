package usecase

import (
	"0700-express-web-api/interface/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthUsecase struct {
	authRepository *repository.AuthRepository
	jwtSecret      string
}

type loginResponse struct {
	UUID         string `json:"uuid"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func CreateAuthUsecase(authRepository *repository.AuthRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		authRepository: authRepository,
		jwtSecret:      jwtSecret,
	}
}

func (usecase *AuthUsecase) Login(email, password string) (*loginResponse, error) {
	user, err := usecase.authRepository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := usecase.generateToken(user.ID.String(), "access", time.Now().Add(time.Hour))
	if err != nil {
		return nil, err
	}

	refreshToken, err := usecase.generateToken(user.ID.String(), "refresh", time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, err
	}

	return &loginResponse{
		UUID:         user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (usecase *AuthUsecase) SignUp(email, password string) error {
	return nil
}

func (usecase *AuthUsecase) generateToken(userId string, tokenType string, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userId,
		"type": tokenType,
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(usecase.jwtSecret))
}
