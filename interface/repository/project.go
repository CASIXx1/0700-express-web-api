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

func (repository *ProjectRepository) FindProjects(ctx context.Context, userID string, limit int, offset int) ([]*ent.Project, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return repository.client.Project.
		Query().
		Limit(limit).
		Offset(offset).
		Where(entProject.UserID(id)).
		Order(entProject.BySortOrder()).
		All(ctx)
}

func (repository *ProjectRepository) FindProjectBySlug(ctx context.Context, userID string, slug string) (*ent.Project, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return repository.client.Project.
		Query().
		Where(entProject.UserID(id)).
		Where(entProject.Slug(slug)).
		Only(ctx)
}

func (repository *ProjectRepository) CreateProject(ctx context.Context, slug string, userID uuid.UUID, sortOrder int) error {
	return repository.client.Project.
		Create().
		SetSlug(slug).
		SetUserID(userID).
		SetSortOrder(sortOrder).
		Exec(ctx)
}

func (repository *ProjectRepository) CountProjects(ctx context.Context, userID string) (int, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	return repository.client.Project.
		Query().
		Where(entProject.UserID(id)).
		Count(ctx)
}
