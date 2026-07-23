package seed

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserSeeder struct{}

func NewUserSeeder() *UserSeeder {
	return &UserSeeder{}
}

func (seeder *UserSeeder) Run(ctx context.Context, client *ent.Client) error {
	userRepository := repository.CreateUserRepository(client)

	exists, err := userRepository.FindUserByEmail(ctx, "test@example.com")
	if exists != nil {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = userRepository.CreateUser(ctx, "test", "test@example.com", string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}
