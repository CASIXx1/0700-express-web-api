package handler

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/middleware"
	"0700-express-web-api/usecase"
	"log"
	"net/http"
	"strconv"
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
	CreatedAt  time.Time            `json:"createdAt"`
	UpdatedAt  time.Time            `json:"updatedAt"`
	Deadline   *time.Time           `json:"deadline"`
	StartingAt *time.Time           `json:"startingAt"`
	StartedAt  *time.Time           `json:"startedAt"`
	FinishedAt *time.Time           `json:"finishedAt"`
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

	limit := 1
	limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
	if err != nil {
		writeResponse(writer, http.StatusBadRequest, errorResponse{
			Message: err.Error(),
		})
		return
	}

	page := 1
	page, err = strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		writeResponse(writer, http.StatusBadRequest, errorResponse{
			Message: err.Error(),
		})
		return
	}

	totalCount, err := handler.projectUsecase.CountProjects(request.Context(), userID)
	if err != nil {
		writeResponse(writer, http.StatusInternalServerError, errorResponse{
			Message: err.Error(),
		})
		return
	}

	offset := (page - 1) * limit
	hasPrevious := page > 1
	hasNext := offset+limit < totalCount

	projects, err := handler.projectUsecase.FindProjects(request.Context(), userID, limit, offset)
	if err != nil {
		log.Printf("failed to find projects: %v", err)

		writeResponse(writer, http.StatusUnauthorized, errorResponse{
			Message: "unauthorized",
		})
		return
	}

	writeResponse(writer, http.StatusOK, paginatedResponse[[]projectResponse]{
		Data: projectResponses(projects),
		PageInfo: paginationResponse{
			Page:        page,
			Limit:       limit,
			TotalCount:  totalCount,
			HasPrevious: hasPrevious,
			HasNext:     hasNext,
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
		CreatedAt:  project.CreatedAt,
		UpdatedAt:  project.UpdatedAt,
		Deadline:   project.Deadline,
		StartingAt: project.StartingAt,
		StartedAt:  project.StartedAt,
		FinishedAt: project.FinishedAt,
	}
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
