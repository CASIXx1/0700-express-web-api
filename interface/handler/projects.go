package handler

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/usecase"
	"log"
	"net/http"
)

type ProjectHandler struct {
	projectUsecase *usecase.ProjectUsecase
}

type projectResponse struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

func NewProjectHandler(projectUsecase *usecase.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{
		projectUsecase: projectUsecase,
	}
}

func (handler *ProjectHandler) FindProjects(writer http.ResponseWriter, request *http.Request) {
	accessToken, err := bearerToken(request)
	if err != nil {
		writeResponse(writer, http.StatusUnauthorized, errorResponse{
			Message: err.Error(),
		})
		return
	}

	projects, err := handler.projectUsecase.FindProjects(request.Context(), accessToken)
	if err != nil {
		log.Printf("failed to find projects: %v", err)

		writeResponse(writer, http.StatusUnauthorized, errorResponse{
			Message: "unauthorized",
		})
		return
	}

	writeResponse(writer, http.StatusOK, normalResponse[[]projectResponse]{
		Data: projectResponses(projects),
	})
}

func projectResponses(projects []*ent.Project) []projectResponse {
	responses := []projectResponse{}

	for _, project := range projects {
		responses = append(responses, projectResponse{
			ID:   project.ID.String(),
			Slug: project.Slug,
		})
	}

	return responses
}
