package repository

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/ent/project"
	"context"
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
