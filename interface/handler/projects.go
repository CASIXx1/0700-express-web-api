package handler

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/middleware"
	"0700-express-web-api/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	projectUsecase *usecase.ProjectUsecase
}

type projectResponse struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
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
		Data: projectResponse{
			ID:   project.ID.String(),
			Slug: project.Slug,
			Name: project.Name,
		},
	})
}

func projectResponses(projects []*ent.Project) []projectResponse {
	responses := []projectResponse{}

	for _, project := range projects {
		responses = append(responses, projectResponse{
			ID:   project.ID.String(),
			Slug: project.Slug,
			Name: project.Name,
		})
	}

	return responses
}
