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
	userID := uuid.New()
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
				repository.EXPECT().CountTasks(ctx, userID, statuses).Return(3, nil)
				repository.EXPECT().FindTasks(ctx, userID, statuses, 1, 1).Return(tasks, nil)
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
				repository.EXPECT().CountTasks(ctx, userID, statuses).Return(0, countTasksError)
			},
			expectedError: countTasksError,
		},
		{
			name: "error case: find tasks failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().CountTasks(ctx, userID, statuses).Return(3, nil)
				repository.EXPECT().FindTasks(ctx, userID, statuses, 1, 1).Return(nil, findTasksError)
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
			result, err := taskUsecase.FindTasks(ctx, userID, statuses, 2, 1)

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
	userID := uuid.New()
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
				repository.EXPECT().CreateTask(ctx, userID, input).Return(task, nil)
			},
			expectedTask: task,
		},
		{
			name: "error case: create task failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().CreateTask(ctx, userID, input).Return(nil, createTaskError)
			},
			expectedError: createTaskError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockTaskRepository(ctrl)
			test.setup(repository)

			result, err := NewTaskUsecase(repository).CreateTask(ctx, userID, input)

			assertTaskResult(t, result, err, test.expectedTask, test.expectedError)
		})
	}
}

func TestTaskUsecaseFindTaskByID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	taskID := uuid.New()
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
				repository.EXPECT().FindTaskByID(ctx, userID, taskID).Return(task, nil)
			},
			expectedTask: task,
		},
		{
			name: "error case: find task failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().FindTaskByID(ctx, userID, taskID).Return(nil, findTaskError)
			},
			expectedError: findTaskError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockTaskRepository(ctrl)
			test.setup(repository)

			result, err := NewTaskUsecase(repository).FindTaskByID(ctx, userID, taskID)

			assertTaskResult(t, result, err, test.expectedTask, test.expectedError)
		})
	}
}

func TestTaskUsecaseUpdateTask(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	taskID := uuid.New()
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
				repository.EXPECT().UpdateTask(ctx, userID, taskID, input).Return(task, nil)
			},
			expectedTask: task,
		},
		{
			name: "error case: update task failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().UpdateTask(ctx, userID, taskID, input).Return(nil, updateTaskError)
			},
			expectedError: updateTaskError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockTaskRepository(ctrl)
			test.setup(repository)

			result, err := NewTaskUsecase(repository).UpdateTask(ctx, userID, taskID, input)

			assertTaskResult(t, result, err, test.expectedTask, test.expectedError)
		})
	}
}

func TestTaskUsecaseDeleteTask(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	taskID := uuid.New()
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
				repository.EXPECT().DeleteTask(ctx, userID, taskID).Return(task, nil)
			},
			expectedTask: task,
		},
		{
			name: "error case: delete task failed",
			setup: func(repository *MockTaskRepository) {
				repository.EXPECT().DeleteTask(ctx, userID, taskID).Return(nil, deleteTaskError)
			},
			expectedError: deleteTaskError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockTaskRepository(ctrl)
			test.setup(repository)

			result, err := NewTaskUsecase(repository).DeleteTask(ctx, userID, taskID)

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
