package usecase

import (
	"0700-express-web-api/ent"
	"context"
	"testing"
)

type fakeProjectRepository struct {
	receivedUserID string
	receivedLimit  int
	receivedOffset int
	totalCount     int
	projects       []*ent.Project
}

func (repository *fakeProjectRepository) FindProjects(ctx context.Context, userID string, limit int, offset int) ([]*ent.Project, error) {
	repository.receivedUserID = userID
	repository.receivedLimit = limit
	repository.receivedOffset = offset

	return repository.projects, nil
}

func (repository *fakeProjectRepository) FindProjectBySlug(ctx context.Context, userID string, slug string) (*ent.Project, error) {
	return nil, nil
}

func (repository *fakeProjectRepository) CountProjects(ctx context.Context, userID string) (int, error) {
	return repository.totalCount, nil
}

func TestProjectUsecaseFindProjects(t *testing.T) {
	repository := &fakeProjectRepository{
		totalCount: 3,
		projects:   []*ent.Project{},
	}
	projectUsecase := NewProjectUsecase(repository)

	result, err := projectUsecase.FindProjects(context.Background(), "user-id", 2, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_ = result
}
