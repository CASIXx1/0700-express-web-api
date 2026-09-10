package usecase

import (
	"0700-express-web-api/ent"
	entTask "0700-express-web-api/ent/task"
	"0700-express-web-api/interface/repository"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTaskUsecaseFindTasks(t *testing.T) {
	ctx := context.Background()
	statuses := []string{"scheduled"}
	tasks := []*ent.Task{}
	countTasksError := errors.New("failed to count tasks")
	findTasksError := errors.New("failed to find tasks")

	tests := []struct {
		name           string
		setup          func(repository *MockTaskRepository)
		expectedResult *TaskListResult
		expectedError  error
	}{
		{
			name: "normal case: find tasks",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().CountTasks(ctx, "user-id", statuses).Return(3, nil)
				repository.EXPECT().FindTasks(ctx, "user-id", statuses, 1, 1).Return(tasks, nil)
			},
			expectedResult: &TaskListResult{
				Tasks: tasks,
				PageInfo: PageInfo{
					TotalCount:  3,
					HasPrevious: true,
					HasNext:     true,
				},
			},
		},
		{
			name: "error case: count tasks failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().CountTasks(ctx, "user-id", statuses).Return(0, countTasksError)
			},
			expectedError: countTasksError,
		},
		{
			name: "error case: find tasks failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().CountTasks(ctx, "user-id", statuses).Return(3, nil)
				repository.EXPECT().FindTasks(ctx, "user-id", statuses, 1, 1).Return(nil, findTasksError)
			},
			expectedError: findTasksError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockTaskRepository(ctrl)
			test.setup(repository)

			taskUsecase := NewTaskUsecase(repository)
			result, err := taskUsecase.FindTasks(ctx, "user-id", statuses, 2, 1)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func TestTaskUsecaseCreateTask(t *testing.T) {
	ctx := context.Background()
	input := repository.CreateTaskInput{
		Title:     "task title",
		Status:    entTask.StatusScheduled,
		ProjectID: uuid.New(),
	}
	task := &ent.Task{}
	createTaskError := errors.New("failed to create task")

	tests := []struct {
		name          string
		setup         func(repository *MockTaskRepository)
		expectedTask  *ent.Task
		expectedError error
	}{
		{
			name: "normal case: create task",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().CreateTask(ctx, "user-id", input).Return(task, nil)
			},
			expectedTask: task,
		},
		{
			name: "error case: create task failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().CreateTask(ctx, "user-id", input).Return(nil, createTaskError)
			},
			expectedError: createTaskError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockTaskRepository(ctrl)
			test.setup(repository)

			result, err := NewTaskUsecase(repository).CreateTask(ctx, "user-id", input)

			assertTaskResult(t, result, err, test.expectedTask, test.expectedError)
		})
	}
}

func TestTaskUsecaseFindTaskByID(t *testing.T) {
	ctx := context.Background()
	task := &ent.Task{}
	findTaskError := errors.New("failed to find task")

	tests := []struct {
		name          string
		setup         func(repository *MockTaskRepository)
		expectedTask  *ent.Task
		expectedError error
	}{
		{
			name: "normal case: find task",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().FindTaskByID(ctx, "user-id", "task-id").Return(task, nil)
			},
			expectedTask: task,
		},
		{
			name: "error case: find task failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().FindTaskByID(ctx, "user-id", "task-id").Return(nil, findTaskError)
			},
			expectedError: findTaskError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockTaskRepository(ctrl)
			test.setup(repository)

			result, err := NewTaskUsecase(repository).FindTaskByID(ctx, "user-id", "task-id")

			assertTaskResult(t, result, err, test.expectedTask, test.expectedError)
		})
	}
}

func TestTaskUsecaseUpdateTask(t *testing.T) {
	ctx := context.Background()
	title := "updated task title"
	input := repository.UpdateTaskInput{Title: &title}
	task := &ent.Task{}
	updateTaskError := errors.New("failed to update task")

	tests := []struct {
		name          string
		setup         func(repository *MockTaskRepository)
		expectedTask  *ent.Task
		expectedError error
	}{
		{
			name: "normal case: update task",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().UpdateTask(ctx, "user-id", "task-id", input).Return(task, nil)
			},
			expectedTask: task,
		},
		{
			name: "error case: update task failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().UpdateTask(ctx, "user-id", "task-id", input).Return(nil, updateTaskError)
			},
			expectedError: updateTaskError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockTaskRepository(ctrl)
			test.setup(repository)

			result, err := NewTaskUsecase(repository).UpdateTask(ctx, "user-id", "task-id", input)

			assertTaskResult(t, result, err, test.expectedTask, test.expectedError)
		})
	}
}

func TestTaskUsecaseDeleteTask(t *testing.T) {
	ctx := context.Background()
	task := &ent.Task{}
	deleteTaskError := errors.New("failed to delete task")

	tests := []struct {
		name          string
		setup         func(repository *MockTaskRepository)
		expectedTask  *ent.Task
		expectedError error
	}{
		{
			name: "normal case: delete task",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().DeleteTask(ctx, "user-id", "task-id").Return(task, nil)
			},
			expectedTask: task,
		},
		{
			name: "error case: delete task failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().DeleteTask(ctx, "user-id", "task-id").Return(nil, deleteTaskError)
			},
			expectedError: deleteTaskError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockTaskRepository(ctrl)
			test.setup(repository)

			result, err := NewTaskUsecase(repository).DeleteTask(ctx, "user-id", "task-id")

			assertTaskResult(t, result, err, test.expectedTask, test.expectedError)
		})
	}
}

func assertTaskResult(t *testing.T, result *ent.Task, err error, expectedTask *ent.Task, expectedError error) {
	t.Helper()

	if expectedError != nil {
		require.ErrorIs(t, err, expectedError)
		assert.Nil(t, result)
		return
	}

	require.NoError(t, err)
	assert.Equal(t, expectedTask, result)
}
