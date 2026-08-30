package handler

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/usecase"
	"log"
	"net/http"
)

type TaskHandler struct {
	taskUsecase *usecase.TaskUsecase
}

type taskResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	FinishedAt  *string  `json:"finishedAt"`
	StartedAt   *string  `json:"startedAt"`
	ArchivedAt  *string  `json:"archivedAt"`
	StartingAt  *string  `json:"startingAt"`
	Deadline    *string  `json:"deadline"`
	Parent      *string  `json:"parent"`
	Children    []string `json:"children"`
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

	status := request.URL.Query().Get("status")
	result, err := handler.taskUsecase.FindTasks(request.Context(), userID, status, paginationRequest.Page, paginationRequest.Limit)
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

func taskResponses(tasks []*ent.Task) []taskResponse {
	responses := []taskResponse{}

	for _, task := range tasks {
		responses = append(responses, taskResponseFromTask(task))
	}

	return responses
}

func taskResponseFromTask(task *ent.Task) taskResponse {
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
		Parent:      nil,
		Children:    []string{},
	}
}
