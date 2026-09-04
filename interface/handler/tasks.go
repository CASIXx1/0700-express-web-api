package handler

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"
	"0700-express-web-api/usecase"
	"errors"
	"log"
	"net/http"

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

func NewTaskHandler(taskUsecase *usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{
		taskUsecase: taskUsecase,
	}
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
