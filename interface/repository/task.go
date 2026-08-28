package repository

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/ent/project"
	"context"

	"github.com/google/uuid"
)

type TaskRepository struct {
	client *ent.Client
}

func NewTaskRepository(client *ent.Client) *TaskRepository {
	return &TaskRepository{
		client: client,
	}
}

func (repository *TaskRepository) CreateTask(ctx context.Context, title string, projectSlug string) error {
	return repository.client.Task.
		Create().
		SetTitle(title).
		SetProject(repository.client.Project.
			Query().
			Where(project.SlugEQ(projectSlug)).
			OnlyX(ctx)).
		Exec(ctx)
}

func (repository *ProjectRepository) FindTasks(ctx context.Context, userID string, limit int, offset int) ([]*ent.Project, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return repository.client.Task.
		Query().
		Limit(limit).
		Offset(offset).
		Where(entTask.UserID(id)).
		Order(entTask.BySortOrder()).
		All(ctx)
}
