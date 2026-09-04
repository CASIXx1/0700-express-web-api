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

func (repository *TaskRepository) FindTasks(ctx context.Context, userID string, statuses []string, limit int, offset int) ([]*ent.Task, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	query := repository.client.Task.
		Query().
		Limit(limit).
		Offset(offset).
		Where(entTask.HasProjectWith(entProject.UserID(id))).
		WithProject().
		Order(entTask.ByID())

	if len(statuses) > 0 {
		query.Where(entTask.StatusIn(taskStatuses(statuses)...))
	}

	return query.All(ctx)
}

func (repository *TaskRepository) CountTasks(ctx context.Context, userID string, statuses []string) (int, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	query := repository.client.Task.
		Query().
		Where(entTask.HasProjectWith(entProject.UserID(id)))

	if len(statuses) > 0 {
		query.Where(entTask.StatusIn(taskStatuses(statuses)...))
	}

	return query.Count(ctx)
}

func (repository *TaskRepository) FindTaskByID(ctx context.Context, userID string, taskID string) (*ent.Task, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return nil, err
	}

	task, err := repository.client.Task.
		Query().
		Where(entTask.ID(taskUUID)).
		Where(entTask.HasProjectWith(entProject.UserID(userUUID))).
		WithProject().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return task, nil
}

func (repository *TaskRepository) DeleteTask(ctx context.Context, userID string, taskID string) (*ent.Task, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return nil, err
	}

	task, err := repository.client.Task.
		Query().
		Where(entTask.ID(taskUUID)).
		Where(entTask.HasProjectWith(entProject.UserID(userUUID))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if err := repository.client.Task.DeleteOne(task).Exec(ctx); err != nil {
		return nil, err
	}

	return task, nil
}

func taskStatuses(statuses []string) []entTask.Status {
	values := []entTask.Status{}
	for _, status := range statuses {
		values = append(values, entTask.Status(status))
	}

	return values
}
