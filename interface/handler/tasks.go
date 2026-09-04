package handler

import (
	"0700-express-web-api/ent"
	entTask "0700-express-web-api/ent/task"
	"0700-express-web-api/interface/repository"
	"0700-express-web-api/usecase"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	taskUsecase *usecase.TaskUsecase
}

type taskResponse struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
	FinishedAt  *string          `json:"finishedAt"`
	StartedAt   *string          `json:"startedAt"`
	ArchivedAt  *string          `json:"archivedAt"`
	StartingAt  *string          `json:"startingAt"`
	Deadline    *string          `json:"deadline"`
	Project     *projectResponse `json:"project,omitempty"`
	Parent      *string          `json:"parent"`
	Children    []string         `json:"children"`
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ProjectID   string `json:"projectId"`
	StartingAt  string `json:"startingAt"`
	Deadline    string `json:"deadline"`
}

type updateTaskRequest struct {
	Title       *string `json:"title"`
	Kind        *string `json:"kind"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	ProjectID   *string `json:"projectId"`
	StartingAt  *string `json:"startingAt"`
	Deadline    *string `json:"deadline"`
}

func NewTaskHandler(taskUsecase *usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{
		taskUsecase: taskUsecase,
	}
}

func (handler *TaskHandler) CreateTask(writer http.ResponseWriter, request *http.Request) {
	userID, ok := userIDFromContext(request.Context())
	if !ok || userID == "" {
		WriteResponse(writer, http.StatusUnauthorized, ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	var body createTaskRequest
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	input, err := createTaskInputFromRequest(body)
	if err != nil {
		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	err = handler.taskUsecase.CreateTask(request.Context(), userID, input)
	if err != nil {
		log.Printf("failed to create task: %v", err)

		statusCode := http.StatusInternalServerError
		if errors.Is(err, repository.ErrNotFound) {
			statusCode = http.StatusBadRequest
		}

		WriteResponse(writer, statusCode, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	writer.WriteHeader(http.StatusCreated)
}

func (handler *TaskHandler) FindTasks(writer http.ResponseWriter, request *http.Request) {
	userID, ok := userIDFromContext(request.Context())
	if !ok || userID == "" {
		WriteResponse(writer, http.StatusUnauthorized, ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	paginationRequest, err := parsePaginationParams(request.URL.Query())
	if err != nil {
		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	statuses := request.URL.Query()["status"]
	result, err := handler.taskUsecase.FindTasks(request.Context(), userID, statuses, paginationRequest.Page, paginationRequest.Limit)
	if err != nil {
		log.Printf("failed to find tasks: %v", err)

		WriteResponse(writer, http.StatusUnauthorized, ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	WriteResponse(writer, http.StatusOK, paginatedResponse[[]taskResponse]{
		Data: taskResponses(result.Tasks),
		PageInfo: paginationResponse{
			Page:        paginationRequest.Page,
			Limit:       paginationRequest.Limit,
			TotalCount:  result.PageInfo.TotalCount,
			HasPrevious: result.PageInfo.HasPrevious,
			HasNext:     result.PageInfo.HasNext,
		},
	})
}

func (handler *TaskHandler) DeleteTask(writer http.ResponseWriter, request *http.Request) {
	userID, ok := userIDFromContext(request.Context())
	if !ok || userID == "" {
		WriteResponse(writer, http.StatusUnauthorized, ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	taskID := mux.Vars(request)["id"]
	if taskID == "" {
		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: "missing task id",
		})
		return
	}

	task, err := handler.taskUsecase.DeleteTask(request.Context(), userID, taskID)
	if err != nil {
		log.Printf("failed to delete task: %v", err)

		if errors.Is(err, repository.ErrNotFound) {
			WriteResponse(writer, http.StatusNotFound, ErrorResponse{
				Message: "task not found",
			})
			return
		}

		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	WriteResponse(writer, http.StatusOK, normalResponse[taskResponse]{
		Data: taskResponseFromTask(task),
	})
}

func (handler *TaskHandler) FindTaskByID(writer http.ResponseWriter, request *http.Request) {
	userID, ok := userIDFromContext(request.Context())
	if !ok || userID == "" {
		WriteResponse(writer, http.StatusUnauthorized, ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	taskID := mux.Vars(request)["id"]
	if taskID == "" {
		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: "missing task id",
		})
		return
	}

	task, err := handler.taskUsecase.FindTaskByID(request.Context(), userID, taskID)
	if err != nil {
		log.Printf("failed to find task: %v", err)

		if errors.Is(err, repository.ErrNotFound) {
			WriteResponse(writer, http.StatusNotFound, ErrorResponse{
				Message: "task not found",
			})
			return
		}

		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	WriteResponse(writer, http.StatusOK, normalResponse[taskResponse]{
		Data: taskResponseFromTask(task),
	})
}

func (handler *TaskHandler) UpdateTask(writer http.ResponseWriter, request *http.Request) {
	userID, ok := userIDFromContext(request.Context())
	if !ok || userID == "" {
		WriteResponse(writer, http.StatusUnauthorized, ErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	taskID := mux.Vars(request)["id"]
	if taskID == "" {
		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: "missing task id",
		})
		return
	}

	var body updateTaskRequest
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	input, err := updateTaskInputFromRequest(body)
	if err != nil {
		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	task, err := handler.taskUsecase.UpdateTask(request.Context(), userID, taskID, input)
	if err != nil {
		log.Printf("failed to update task: %v", err)

		if errors.Is(err, repository.ErrNotFound) {
			WriteResponse(writer, http.StatusNotFound, ErrorResponse{
				Message: "task not found",
			})
			return
		}

		WriteResponse(writer, http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	WriteResponse(writer, http.StatusOK, normalResponse[taskResponse]{
		Data: taskResponseFromTask(task),
	})
}

func taskResponses(tasks []*ent.Task) []taskResponse {
	responses := []taskResponse{}

	for _, task := range tasks {
		responses = append(responses, taskResponseFromTask(task))
	}

	return responses
}

func taskResponseFromTask(task *ent.Task) taskResponse {
	var project *projectResponse
	if task.Edges.Project != nil {
		response := projectResponseFromProject(task.Edges.Project)
		project = &response
	}

	return taskResponse{
		ID:          task.ID.String(),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status.String(),
		CreatedAt:   formatDateTime(task.CreatedAt),
		UpdatedAt:   formatDateTime(task.UpdatedAt),
		FinishedAt:  formatOptionalDateTime(task.FinishedAt),
		StartedAt:   formatOptionalDateTime(task.StartedAt),
		ArchivedAt:  formatOptionalDateTime(task.ArchivedAt),
		StartingAt:  formatOptionalDateTime(task.StartingAt),
		Deadline:    formatOptionalDateTime(task.Deadline),
		Project:     project,
		Parent:      nil,
		Children:    []string{},
	}
}

func createTaskInputFromRequest(request createTaskRequest) (repository.CreateTaskInput, error) {
	if request.Title == "" {
		return repository.CreateTaskInput{}, errors.New("missing title")
	}

	if request.Kind != "task" {
		return repository.CreateTaskInput{}, errors.New("invalid kind")
	}

	status := entTask.Status(request.Status)
	err := entTask.StatusValidator(status)
	if err != nil {
		return repository.CreateTaskInput{}, err
	}

	projectID, err := uuid.Parse(request.ProjectID)
	if err != nil {
		return repository.CreateTaskInput{}, err
	}

	startedAt, err := time.Parse("2026-01-01", request.StartingAt)
	if err != nil {
		return repository.CreateTaskInput{}, err
	}

	deadline, err := time.Parse("2026-01-01", request.Deadline)
	if err != nil {
		return repository.CreateTaskInput{}, err
	}

	return repository.CreateTaskInput{
		Title:       request.Title,
		Description: request.Description,
		Status:      status,
		ProjectID:   projectID,
		StartingAt:  &startedAt,
		Deadline:    &deadline,
	}, nil
}

func updateTaskInputFromRequest(request updateTaskRequest) (repository.UpdateTaskInput, error) {
	return repository.UpdateTaskInput{
		Title:       request.Title,
		Description: request.Description,
		Status:      request.Status,
		ProjectID:   request.ProjectID,
		StartingAt:  request.StartingAt,
		Deadline:    request.Deadline,
	}, nil
}
