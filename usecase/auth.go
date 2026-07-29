package usecase

import (
	"0700-express-web-api/ent"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserAlreadyExists = errors.New("user already exists")

type AuthUsecase struct {
	userRepository   UserRepository
	passwordVerifier passwordVerifier
	tokenGenerator   tokenGenerator
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

type passwordVerifier interface {
	Verify(userPassword, requestPassword string) error
}

type tokenGenerator interface {
	GenerateAccessToken(userIdentifier string) (string, error)
	GenerateRefreshToken(userIdentifier string) (string, error)
}

func NewAuthUsecase(userRepository UserRepository, passwordVerifier passwordVerifier, tokenGenerator tokenGenerator) *AuthUsecase {
	return &AuthUsecase{
		userRepository:   userRepository,
		passwordVerifier: passwordVerifier,
		tokenGenerator:   tokenGenerator,
	}
}

func (usecase *AuthUsecase) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	user, err := usecase.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := usecase.passwordVerifier.Verify(user.Password, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := usecase.tokenGenerator.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := usecase.tokenGenerator.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
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

	accessToken, err := usecase.tokenGenerator.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := usecase.tokenGenerator.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		UUID:         user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateToken(userId string, tokenType string, expiresAt time.Time, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userId,
		"type": tokenType,
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
