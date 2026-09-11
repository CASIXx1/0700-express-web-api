package usecase

import (
	"0700-express-web-api/ent"
	"context"

	"github.com/google/uuid"
)

type ProjectUsecase struct {
	projectRepository ProjectRepository
}

type ProjectListResult struct {
	Projects []*ent.Project
	PageInfo PageInfo
}

type PageInfo struct {
	TotalCount  int
	HasPrevious bool
	HasNext     bool
}

type ProjectRepository interface {
	FindProjects(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*ent.Project, error)
	FindProjectBySlug(ctx context.Context, userID uuid.UUID, slug string) (*ent.Project, error)
	CountProjects(ctx context.Context, userID uuid.UUID) (int, error)
}

func NewProjectUsecase(projectRepository ProjectRepository) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepository: projectRepository,
	}
}

func (usecase *ProjectUsecase) FindProjects(ctx context.Context, userID uuid.UUID, page int, limit int) (*ProjectListResult, error) {
	offset := (page - 1) * limit

	totalCount, err := usecase.projectRepository.CountProjects(ctx, userID)
	if err != nil {
		return nil, err
	}

	projects, err := usecase.projectRepository.FindProjects(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &ProjectListResult{
		Projects: projects,
		PageInfo: PageInfo{
			TotalCount:  totalCount,
			HasPrevious: page > 1,
			HasNext:     offset+limit < totalCount,
		},
	}, nil
}

func (usecase *ProjectUsecase) FindProjectBySlug(ctx context.Context, userID uuid.UUID, slug string) (*ent.Project, error) {
	return usecase.projectRepository.FindProjectBySlug(ctx, userID, slug)
}
