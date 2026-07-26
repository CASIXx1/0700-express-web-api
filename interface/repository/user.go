package repository

import (
	"0700-express-web-api/ent"
	entuser "0700-express-web-api/ent/user"
	"context"
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

func (repository *UserRepository) CreateUser(ctx context.Context, username, email, hashedPassword string) (*ent.User, error) {
	return repository.dbClient.User.
		Create().
		SetUsername(username).
		SetEmail(email).
		SetPassword(hashedPassword).
		Save(ctx)
}
