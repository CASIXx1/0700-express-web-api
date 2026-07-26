package usecase

import (
	"0700-express-web-api/ent"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserAlreadyExists = errors.New("user already exists")

type AuthUsecase struct {
	userRepository UserRepository
	jwtSecret      string
}

type AuthResult struct {
	UUID         string
	AccessToken  string
	RefreshToken string
}

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*ent.User, error)
	CreateUser(ctx context.Context, username, email, hashedPassword string) (*ent.User, error)
}

func NewAuthUsecase(userRepository UserRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (usecase *AuthUsecase) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	user, err := usecase.userRepository.FindUserByEmail(ctx, email)
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

	return &AuthResult{
		UUID:         user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (usecase *AuthUsecase) SignUp(ctx context.Context, username, email, password string) (*AuthResult, error) {
	user, err := usecase.userRepository.FindUserByEmail(ctx, email)
	if user != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err = usecase.userRepository.CreateUser(ctx, username, email, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	accessToken, err := usecase.generateToken(user.ID.String(), "access", time.Now().Add(time.Hour))
	if err != nil {
		return nil, err
	}

	refreshToken, err := usecase.generateToken(user.ID.String(), "refresh", time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		UUID:         user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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
