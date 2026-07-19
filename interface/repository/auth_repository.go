package repository

import (
	"0700-express-web-api/ent"
	entuser "0700-express-web-api/ent/user"
	"context"
)

type AuthRepository struct {
	dbClient *ent.Client
}

func CreateAuthRepository(dbClient *ent.Client) *AuthRepository {
	return &AuthRepository{
		dbClient: dbClient,
	}
}

func (repository *AuthRepository) FindUserByEmail(email string) (*ent.User, error) {
	return repository.dbClient.User.
		Query().
		Where(entuser.EmailEQ(email)).
		First(context.Background())
}
