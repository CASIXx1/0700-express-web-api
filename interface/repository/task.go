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

func (repository *TaskRepository) CreateTask(ctx context.Context, userID string, input CreateTaskInput) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	project, err := repository.client.Project.
		Query().
		Where(entProject.ID(input.ProjectID)).
		Where(entProject.UserID(userUUID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}

		return err
	}

	return repository.client.Task.
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

func (repository *TaskRepository) UpdateTask(ctx context.Context, userID string, taskID string, input UpdateTaskInput) (*ent.Task, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return nil, err
	}

	currentTask, err := repository.FindTaskByID(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	project := currentTask.Edges.Project
	if input.ProjectID != nil {
		project, err = repository.client.Project.
			Query().
			Where(entProject.ID(*input.ProjectID)).
			Where(entProject.UserID(userUUID)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, ErrNotFound
			}

			return nil, err
		}
	}

	update := repository.client.Task.
		UpdateOneID(taskUUID).
		Where(entTask.HasProjectWith(entProject.UserID(userUUID)))

	if input.Title != nil {
		update.SetTitle(*input.Title)
	}
	if input.Description != nil {
		update.SetDescription(*input.Description)
	}
	if input.Status != nil {
		update.SetStatus(*input.Status)
	}
	if input.FinishedAt != nil {
		update.SetFinishedAt(*input.FinishedAt)
	}
	if input.StartedAt != nil {
		update.SetStartedAt(*input.StartedAt)
	}
	if input.ArchivedAt != nil {
		update.SetArchivedAt(*input.ArchivedAt)
	}
	if input.StartingAt != nil {
		update.SetStartingAt(*input.StartingAt)
	}
	if input.Deadline != nil {
		update.SetDeadline(*input.Deadline)
	}
	if input.ProjectID != nil {
		update.SetProjectID(project.ID)
	}

	task, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	task.Edges.Project = project

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
