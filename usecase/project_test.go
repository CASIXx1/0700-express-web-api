package usecase

import (
	"0700-express-web-api/ent"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestProjectUsecaseFindProjects(t *testing.T) {
	ctx := context.Background()
	projects := []*ent.Project{}

	ctrl := gomock.NewController(t)
	repository := NewMockProjectRepository(ctrl)
	repository.EXPECT().CountProjects(ctx, "user-id").Return(3, nil)
	repository.EXPECT().FindProjects(ctx, "user-id", 1, 1).Return(projects, nil)

	projectUsecase := NewProjectUsecase(repository)

	result, err := projectUsecase.FindProjects(ctx, "user-id", 2, 1)
	require.NoError(t, err)

	assert.Equal(t, &ProjectListResult{
		Projects: projects,
		PageInfo: PageInfo{
			TotalCount:  3,
			HasPrevious: true,
			HasNext:     true,
		},
	}, result)
}
