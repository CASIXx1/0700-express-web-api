package usecase

import (
	"0700-express-web-api/ent"
	"context"
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
	Limit       int
	Page        int
	HasPrevious bool
	HasNext     bool
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

func (usecase *ProjectUsecase) FindProjects(ctx context.Context, userID string, page int, limit int) (*ProjectListResult, error) {
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
			Limit:       limit,
			Page:        page,
			HasPrevious: page > 1,
			HasNext:     offset+limit < totalCount,
		},
	}, nil
}

func (usecase *ProjectUsecase) FindProjectBySlug(ctx context.Context, userID string, slug string) (*ent.Project, error) {
	return usecase.projectRepository.FindProjectBySlug(ctx, userID, slug)
}
