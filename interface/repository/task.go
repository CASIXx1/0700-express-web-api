package repository

import (
	"0700-express-web-api/ent"
	entProject "0700-express-web-api/ent/project"
	entTask "0700-express-web-api/ent/task"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateTaskInput struct {
	Title       string
	Description string
	Status      string
	ProjectID   uuid.UUID
	FinishedAt  *time.Time
	StartedAt   *time.Time
	ArchivedAt  *time.Time
	StartingAt  *time.Time
	Deadline    *time.Time
}

type TaskRepository struct {
	client *ent.Client
}

func NewTaskRepository(client *ent.Client) *TaskRepository {
	return &TaskRepository{
		client: client,
	}
}

func (repository *TaskRepository) CreateTask(ctx context.Context, input CreateTaskInput) error {
	return repository.client.Task.
		Create().
		SetTitle(input.Title).
		SetDescription(input.Description).
		SetStatus(entTask.Status(input.Status)).
		SetNillableFinishedAt(input.FinishedAt).
		SetNillableStartedAt(input.StartedAt).
		SetNillableArchivedAt(input.ArchivedAt).
		SetNillableStartingAt(input.StartingAt).
		SetNillableDeadline(input.Deadline).
		SetProjectID(input.ProjectID).
		Exec(ctx)
}

func (repository *TaskRepository) FindTasks(ctx context.Context, userID string, status string, limit int, offset int) ([]*ent.Task, error) {
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
