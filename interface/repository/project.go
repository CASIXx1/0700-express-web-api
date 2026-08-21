package repository

import (
	"0700-express-web-api/ent"
	entProject "0700-express-web-api/ent/project"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateProjectInput struct {
	Name       string
	Slug       string
	Goal       string
	Shouldbe   string
	Color      string
	Deadline   time.Time
	StartingAt time.Time
	StartedAt  time.Time
	FinishedAt time.Time
	UserID     uuid.UUID
	SortOrder  int
}

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

func (repository *ProjectRepository) CreateProject(ctx context.Context, input CreateProjectInput) error {
	return repository.client.Project.
		Create().
		SetName(input.Name).
		SetSlug(input.Slug).
		SetGoal(input.Goal).
		SetShouldbe(input.Shouldbe).
		SetColor(input.Color).
		SetDeadline(input.Deadline).
		SetStartingAt(input.StartingAt).
		SetStartedAt(input.StartedAt).
		SetFinishedAt(input.FinishedAt).
		SetUserID(input.UserID).
		SetSortOrder(input.SortOrder).
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
