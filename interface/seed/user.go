package seed

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/ent/user"
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserSeeder struct{}

func NewUserSeeder() *UserSeeder {
	return &UserSeeder{}
}

func (seeder *UserSeeder) Run(ctx context.Context, client *ent.Client) error {
	exists, err := client.User.
		Query().
		Where(user.EmailEQ("admin@example.com")).
		Exist(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return client.User.
		Create().
		SetUsername("test").
		SetEmail("test@example.com").
		SetPassword(string(hashedPassword)).
		Exec(ctx)
}
