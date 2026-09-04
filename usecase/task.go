package usecase

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"
	"context"
)

type TaskUsecase struct {
	taskRepository TaskRepository
}

type TaskListResult struct {
	Tasks    []*ent.Task
	PageInfo PageInfo
}

type TaskRepository interface {
	CreateTask(ctx context.Context, userID string, input repository.CreateTaskInput) error
	FindTasks(ctx context.Context, userID string, statuses []string, limit int, offset int) ([]*ent.Task, error)
	CountTasks(ctx context.Context, userID string, statuses []string) (int, error)
	FindTaskByID(ctx context.Context, userID string, taskID string) (*ent.Task, error)
	DeleteTask(ctx context.Context, userID string, taskID string) (*ent.Task, error)
}

func NewTaskUsecase(taskRepository TaskRepository) *TaskUsecase {
	return &TaskUsecase{
		taskRepository: taskRepository,
	}
}

func (usecase *TaskUsecase) CreateTask(ctx context.Context, userID string, input repository.CreateTaskInput) error {
	return usecase.taskRepository.CreateTask(ctx, userID, input)
}

func (usecase *TaskUsecase) FindTasks(ctx context.Context, userID string, statuses []string, page int, limit int) (*TaskListResult, error) {
	offset := (page - 1) * limit

	totalCount, err := usecase.taskRepository.CountTasks(ctx, userID, statuses)
	if err != nil {
		return nil, err
	}

	tasks, err := usecase.taskRepository.FindTasks(ctx, userID, statuses, limit, offset)
	if err != nil {
		return nil, err
	}

	return &TaskListResult{
		Tasks: tasks,
		PageInfo: PageInfo{
			TotalCount:  totalCount,
			HasPrevious: page > 1,
			HasNext:     offset+limit < totalCount,
		},
	}, nil
}

func (usecase *TaskUsecase) FindTaskByID(ctx context.Context, userID string, taskID string) (*ent.Task, error) {
	return usecase.taskRepository.FindTaskByID(ctx, userID, taskID)
}

func (usecase *TaskUsecase) DeleteTask(ctx context.Context, userID string, taskID string) (*ent.Task, error) {
	return usecase.taskRepository.DeleteTask(ctx, userID, taskID)
}
