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
	FindProjects(ctx context.Context, userId string, limit int, offset int) ([]*ent.Project, error)
	FindProjectBySlug(ctx context.Context, userId string, slug string) (*ent.Project, error)
	CountProjects(ctx context.Context, userID string) (int, error)
}

func NewProjectUsecase(projectRepository ProjectRepository, tokenVerifier tokenVerifier) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepository: projectRepository,
		tokenVerifier:     tokenVerifier,
	}
}

func (usecase *ProjectUsecase) FindProjects(ctx context.Context, accessToken string, limit int, offset int) ([]*ent.Project, error) {
	userID, err := usecase.tokenVerifier.VerifyAccessToken(accessToken)

	if err != nil {
		return nil, err
	}

	return usecase.projectRepository.FindProjects(ctx, userID, limit, offset)
}

func (usecase *ProjectUsecase) FindProjectBySlug(ctx context.Context, accessToken string, slug string) (*ent.Project, error) {
	userID, err := usecase.tokenVerifier.VerifyAccessToken(accessToken)

	if err != nil {
		return nil, err
	}

	return usecase.projectRepository.FindProjectBySlug(ctx, userID, slug)
}

func (usecase *ProjectUsecase) CountProjects(ctx context.Context, accessToken string) (int, error) {
	userID, err := usecase.tokenVerifier.VerifyAccessToken(accessToken)

	if err != nil {
		return 0, err
	}

	return usecase.projectRepository.CountProjects(ctx, userID)
}
