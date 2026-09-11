package usecase

import (
	"0700-express-web-api/ent"
	"context"

	"github.com/google/uuid"
)

type UserUsecase struct {
	userRepository UserRepository
}

type UserRepository interface {
	FindUserByID(ctx context.Context, userID uuid.UUID) (*ent.User, error)
}

func NewUserUsecase(userRepository UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (usecase *UserUsecase) Me(ctx context.Context, userID uuid.UUID) (*ent.User, error) {
	return usecase.userRepository.FindUserByID(ctx, userID)
}
