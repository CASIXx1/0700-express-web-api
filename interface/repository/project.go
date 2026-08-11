package repository

import (
	"0700-express-web-api/ent"
	entProject "0700-express-web-api/ent/project"
	"context"

	"github.com/google/uuid"
)

type ProjectRepository struct {
	client *ent.Client
}

func NewProjectRepository(client *ent.Client) *ProjectRepository {
	return &ProjectRepository{client}
}

func (repository *ProjectRepository) FindProjects(ctx context.Context, userID string) ([]*ent.Project, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return repository.client.Project.
		Query().
		Where(entProject.UserID(id)).
		All(ctx)
}

func (repository *ProjectRepository) CreateProject(ctx context.Context, slug string, userID uuid.UUID) error {
	return repository.client.Project.
		Create().
		SetSlug(slug).
		SetUserID(userID).
		Exec(ctx)
}
