package usecase

import (
	"0700-express-web-api/ent"
	"context"
)

type ProjectUsecase struct {
	projectRepository ProjectRepository
	tokenVerifier     tokenVerifier
}

type ProjectRepository interface {
	FindProjects(ctx context.Context, userId string, limit int) ([]*ent.Project, error)
}

func NewProjectUsecase(projectRepository ProjectRepository, tokenVerifier tokenVerifier) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepository: projectRepository,
		tokenVerifier:     tokenVerifier,
	}
}

func (usecase *ProjectUsecase) FindProjects(ctx context.Context, accessToken string, limit int) ([]*ent.Project, error) {
	userID, err := usecase.tokenVerifier.VerifyAccessToken(accessToken)

	if err != nil {
		return nil, err
	}

	return usecase.projectRepository.FindProjects(ctx, userID, limit)
}
