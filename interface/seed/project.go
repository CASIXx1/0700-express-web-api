package seed

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"
	"context"
)

type ProjectSeeder struct{}

func NewProjectSeeder() *ProjectSeeder {
	return &ProjectSeeder{}
}

func (seeder *ProjectSeeder) Run(ctx context.Context, client *ent.Client) error {
	projectRepository := repository.CreateProjectRepository(client)

	if err := projectRepository.CreateProject(ctx, "programming"); err != nil {
		return err
	}

	if err := projectRepository.CreateProject(ctx, "design"); err != nil {
		return err
	}

	if err := projectRepository.CreateProject(ctx, "english"); err != nil {
		return err
	}

	return nil
}
