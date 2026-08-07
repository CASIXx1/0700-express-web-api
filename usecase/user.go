package usecase

import (
	"0700-express-web-api/ent"
	"context"
)

type UserUsecase struct {
	userRepository UserRepository
	tokenVerifier  tokenVerifier
}

type UserRepository interface {
	FindUserByID(ctx context.Context, userID string) (*ent.User, error)
}

type tokenVerifier interface {
	VerifyAccessToken(accessToken string) (string, error)
}

func NewUserUsecase(userRepository UserRepository, tokenVerifier tokenVerifier) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
		tokenVerifier:  tokenVerifier,
	}
}

func (usecase *UserUsecase) Me(ctx context.Context, accessToken string) (*ent.User, error) {
	userID, err := usecase.tokenVerifier.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	return usecase.userRepository.FindUserByID(ctx, userID)
}
