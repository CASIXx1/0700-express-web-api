package seed

import (
	"context"

	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"
)

type projectSeeder struct{}

func newProjectSeeder() seeder {
	return &projectSeeder{}
}

func (seeder *projectSeeder) Run(ctx context.Context, client *ent.Client) error {
	projectRepository := repository.NewProjectRepository(client)
	userRepository := repository.NewUserRepository(client)

	user, err := userRepository.FindUserByEmail(ctx, "test@example.com")
	if err != nil {
		return err
	}

	if err := projectRepository.CreateProject(ctx, "programming", user.ID, 1); err != nil {
		return err
	}

	if err := projectRepository.CreateProject(ctx, "english", user.ID, 2); err != nil {
		return err
	}

	if err := projectRepository.CreateProject(ctx, "design", user.ID, 3); err != nil {
		return err
	}

	return nil
}
