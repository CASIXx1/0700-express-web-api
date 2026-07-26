package repository

import (
	"0700-express-web-api/ent"
	"context"
)

type ProjectRepository struct {
	client *ent.Client
}

func NewProjectRepository(client *ent.Client) *ProjectRepository {
	return &ProjectRepository{client}
}

func (repository *ProjectRepository) CreateProject(ctx context.Context, slug string) error {
	return repository.client.Project.
		Create().
		SetSlug(slug).
		Exec(ctx)
}
