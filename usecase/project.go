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
	FindProjects(ctx context.Context) ([]*ent.Project, error)
}

func NewProjectUsecase(projectRepository ProjectRepository, tokenVerifier tokenVerifier) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepository: projectRepository,
		tokenVerifier:     tokenVerifier,
	}
}

// Userが持っているProjectを取得する
func (usecase *ProjectUsecase) FindProjects(ctx context.Context, accessToken string) ([]*ent.Project, error) {
	if _, err := usecase.tokenVerifier.VerifyAccessToken(accessToken); err != nil {
		return nil, err
	}

	return usecase.projectRepository.FindProjects(ctx)
}
