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
	Status      entTask.Status
	ProjectID   uuid.UUID
	FinishedAt  *time.Time
	StartedAt   *time.Time
	ArchivedAt  *time.Time
	StartingAt  *time.Time
	Deadline    *time.Time
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *entTask.Status
	ProjectID   *uuid.UUID
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

func (repository *TaskRepository) CreateTask(ctx context.Context, userID uuid.UUID, input CreateTaskInput) (*ent.Task, error) {
	tx, err := repository.client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	project, err := tx.Project.
		Query().
		Where(entProject.ID(input.ProjectID)).
		Where(entProject.UserID(userID)).
		Only(ctx)
	if err != nil {
		tx.Rollback()

		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	task, err := tx.Task.
		Create().
		SetTitle(input.Title).
		SetDescription(input.Description).
		SetStatus(input.Status).
		SetNillableFinishedAt(input.FinishedAt).
		SetNillableStartedAt(input.StartedAt).
		SetNillableArchivedAt(input.ArchivedAt).
		SetNillableStartingAt(input.StartingAt).
		SetNillableDeadline(input.Deadline).
		SetProjectID(project.ID).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	task = task.Unwrap()
	project = project.Unwrap()
	task.Edges.Project = project

	return task, nil
}

func (repository *TaskRepository) FindTasks(ctx context.Context, userID uuid.UUID, statuses []string, limit int, offset int) ([]*ent.Task, error) {
	query := repository.client.Task.
		Query().
		Limit(limit).
		Offset(offset).
		Where(entTask.HasProjectWith(entProject.UserID(userID))).
		WithProject().
		Order(entTask.ByID())

	if len(statuses) > 0 {
		query.Where(entTask.StatusIn(taskStatuses(statuses)...))
	}

	return query.All(ctx)
}

func (repository *TaskRepository) CountTasks(ctx context.Context, userID uuid.UUID, statuses []string) (int, error) {
	query := repository.client.Task.
		Query().
		Where(entTask.HasProjectWith(entProject.UserID(userID)))

	if len(statuses) > 0 {
		query.Where(entTask.StatusIn(taskStatuses(statuses)...))
	}

	return query.Count(ctx)
}

func (repository *TaskRepository) FindTaskByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*ent.Task, error) {
	task, err := repository.client.Task.
		Query().
		Where(entTask.ID(taskID)).
		Where(entTask.HasProjectWith(entProject.UserID(userID))).
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

func (repository *TaskRepository) UpdateTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, input UpdateTaskInput) (*ent.Task, error) {
	currentTask, err := repository.FindTaskByID(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	project := currentTask.Edges.Project
	if input.ProjectID != nil {
		project, err = repository.client.Project.
			Query().
			Where(entProject.ID(*input.ProjectID)).
			Where(entProject.UserID(userID)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, ErrBadRequest
			}

			return nil, err
		}
	}

	task, err := repository.client.Task.
		UpdateOneID(taskID).
		Where(entTask.HasProjectWith(entProject.UserID(userID))).
		SetNillableTitle(input.Title).
		SetNillableDescription(input.Description).
		SetNillableStatus(input.Status).
		SetNillableFinishedAt(input.FinishedAt).
		SetNillableStartedAt(input.StartedAt).
		SetNillableArchivedAt(input.ArchivedAt).
		SetNillableStartingAt(input.StartingAt).
		SetNillableDeadline(input.Deadline).
		SetNillableProjectID(input.ProjectID).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	task.Edges.Project = project

	return task, nil
}

func (repository *TaskRepository) DeleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*ent.Task, error) {
	task, err := repository.client.Task.
		Query().
		Where(entTask.ID(taskID)).
		Where(entTask.HasProjectWith(entProject.UserID(userID))).
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
