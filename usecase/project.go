package usecase

import (
	"0700-express-web-api/ent"
	"context"
)

type ProjectUsecase struct {
	projectRepository ProjectRepository
}

type ProjectRepository interface {
	FindProjects(ctx context.Context, userID string, limit int, offset int) ([]*ent.Project, error)
	FindProjectBySlug(ctx context.Context, userID string, slug string) (*ent.Project, error)
	CountProjects(ctx context.Context, userID string) (int, error)
}

func NewProjectUsecase(projectRepository ProjectRepository) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepository: projectRepository,
	}
}

func (usecase *ProjectUsecase) FindProjects(ctx context.Context, userID string, limit int, offset int) ([]*ent.Project, error) {
	return usecase.projectRepository.FindProjects(ctx, userID, limit, offset)
}

func (usecase *ProjectUsecase) FindProjectBySlug(ctx context.Context, userID string, slug string) (*ent.Project, error) {
	return usecase.projectRepository.FindProjectBySlug(ctx, userID, slug)
}

func (usecase *ProjectUsecase) CountProjects(ctx context.Context, userID string) (int, error) {
	return usecase.projectRepository.CountProjects(ctx, userID)
}
