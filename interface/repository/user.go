package repository

import (
	"0700-express-web-api/ent"
	entuser "0700-express-web-api/ent/user"
	"context"

	"github.com/google/uuid"
)

type UserRepository struct {
	dbClient *ent.Client
}

func NewUserRepository(dbClient *ent.Client) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}

func (repository *UserRepository) FindUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	return repository.dbClient.User.
		Query().
		Where(entuser.EmailEQ(email)).
		First(ctx)
}

func (repository *UserRepository) FindUserByID(ctx context.Context, userID string) (*ent.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return repository.dbClient.User.Get(ctx, id)
}

func (repository *UserRepository) CreateUser(ctx context.Context, username, email, hashedPassword string) (*ent.User, error) {
	return repository.dbClient.User.
		Create().
		SetUsername(username).
		SetEmail(email).
		SetPassword(hashedPassword).
		Save(ctx)
}
