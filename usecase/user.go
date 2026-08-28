package usecase

import (
	"0700-express-web-api/ent"
	"context"
)

type UserUsecase struct {
	userRepository UserRepository
}

type UserRepository interface {
	FindUserByID(ctx context.Context, userID string) (*ent.User, error)
}

func NewUserUsecase(userRepository UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (usecase *UserUsecase) Me(ctx context.Context, userID string) (*ent.User, error) {
	return usecase.userRepository.FindUserByID(ctx, userID)
}
