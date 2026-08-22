package handler

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/middleware"
	"0700-express-web-api/usecase"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	projectUsecase *usecase.ProjectUsecase
}

type projectResponse struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Slug       string               `json:"slug"`
	Goal       *string              `json:"goal"`
	Shouldbe   *string              `json:"shouldbe"`
	Color      *string              `json:"color"`
	Stats      projectStatsResponse `json:"stats"`
	CreatedAt  string               `json:"createdAt"`
	UpdatedAt  string               `json:"updatedAt"`
	Deadline   *string              `json:"deadline"`
	StartingAt *string              `json:"startingAt"`
	StartedAt  *string              `json:"startedAt"`
	FinishedAt *string              `json:"finishedAt"`
}

type projectStatsResponse struct {
	Total  int                        `json:"total"`
	Kinds  projectStatsKindsResponse  `json:"kinds"`
	States projectStatsStatesResponse `json:"states"`
}

type projectStatsKindsResponse struct {
	Milestone int `json:"milestone"`
	Task      int `json:"task"`
	Total     int `json:"total"`
}

type projectStatsStatesResponse struct {
	Scheduled int `json:"scheduled"`
	Completed int `json:"completed"`
	Archived  int `json:"archived"`
}

func NewProjectHandler(projectUsecase *usecase.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{
		projectUsecase: projectUsecase,
	}
}

func (handler *ProjectHandler) FindProjects(writer http.ResponseWriter, request *http.Request) {
	userID, ok := middleware.UserIDFromContext(request.Context())
	if !ok || userID == "" {
		writeResponse(writer, http.StatusUnauthorized, errorResponse{
			Message: "unauthorized",
		})
		return
	}

	paginationRequest, err := parsePaginationParams(request.URL.Query())
	if err != nil {
		writeResponse(writer, http.StatusBadRequest, errorResponse{
			Message: err.Error(),
		})
		return
	}

	result, err := handler.projectUsecase.FindProjects(request.Context(), userID, paginationRequest.Page, paginationRequest.Limit)
	if err != nil {
		log.Printf("failed to find projects: %v", err)

		writeResponse(writer, http.StatusUnauthorized, errorResponse{
			Message: "unauthorized",
		})
		return
	}

	writeResponse(writer, http.StatusOK, paginatedResponse[[]projectResponse]{
		Data: projectResponses(result.Projects),
		PageInfo: paginationResponse{
			Page:        paginationRequest.Page,
			Limit:       paginationRequest.Limit,
			TotalCount:  result.PageInfo.TotalCount,
			HasPrevious: result.PageInfo.HasPrevious,
			HasNext:     result.PageInfo.HasNext,
		},
	})
}

func (handler *ProjectHandler) FindProjectBySlug(writer http.ResponseWriter, request *http.Request) {
	userID, ok := middleware.UserIDFromContext(request.Context())
	if !ok || userID == "" {
		writeResponse(writer, http.StatusUnauthorized, errorResponse{
			Message: "unauthorized",
		})
		return
	}

	slug := mux.Vars(request)["slug"]
	if slug == "" {
		writeResponse(writer, http.StatusBadRequest, errorResponse{
			Message: "missing slug",
		})
		return
	}

	project, err := handler.projectUsecase.FindProjectBySlug(request.Context(), userID, slug)
	if err != nil {
		log.Printf("failed to find projects: %v", err)

		writeResponse(writer, http.StatusNotFound, errorResponse{
			Message: "project not found",
		})
		return
	}

	writeResponse(writer, http.StatusOK, normalResponse[projectResponse]{
		Data: projectResponseFromProject(project),
	})
}

func projectResponses(projects []*ent.Project) []projectResponse {
	responses := []projectResponse{}

	for _, project := range projects {
		responses = append(responses, projectResponseFromProject(project))
	}

	return responses
}

func projectResponseFromProject(project *ent.Project) projectResponse {
	return projectResponse{
		ID:         project.ID.String(),
		Name:       project.Name,
		Slug:       project.Slug,
		Goal:       project.Goal,
		Shouldbe:   project.Shouldbe,
		Color:      project.Color,
		Stats:      defaultProjectStatsResponse(),
		CreatedAt:  formatDateTime(project.CreatedAt),
		UpdatedAt:  formatDateTime(project.UpdatedAt),
		Deadline:   formatOptionalDateTime(project.Deadline),
		StartingAt: formatOptionalDateTime(project.StartingAt),
		StartedAt:  formatOptionalDateTime(project.StartedAt),
		FinishedAt: formatOptionalDateTime(project.FinishedAt),
	}
}

func formatDateTime(value time.Time) string {
	return value.UTC().Format("2006-01-02T15:04:05.000Z")
}

func formatOptionalDateTime(value *time.Time) *string {
	if value == nil {
		return nil
	}

	formatted := formatDateTime(*value)
	return &formatted
}

func defaultProjectStatsResponse() projectStatsResponse {
	return projectStatsResponse{
		Total: 100,
		Kinds: projectStatsKindsResponse{
			Milestone: 20,
			Task:      30,
			Total:     50,
		},
		States: projectStatsStatesResponse{
			Scheduled: 20,
			Completed: 30,
			Archived:  50,
		},
	}
}
