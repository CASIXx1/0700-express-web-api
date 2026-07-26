package main

import (
	"context"

	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserSeeder struct{}

func NewUserSeeder() *UserSeeder {
	return &UserSeeder{}
}

func (seeder *UserSeeder) Run(ctx context.Context, client *ent.Client) error {
	userRepository := repository.NewUserRepository(client)

	exists, err := userRepository.FindUserByEmail(ctx, "test@example.com")
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if exists != nil {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = userRepository.CreateUser(ctx, "Test User", "test@example.com", string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}
