package usecase

import (
	"0700-express-web-api/ent"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestProjectUsecaseFindProjects(t *testing.T) {
	ctx := context.Background()
	projects := []*ent.Project{}
	countProjectsError := errors.New("failed to count projects")
	findProjectsError := errors.New("failed to find projects")

	tests := []struct {
		name           string
		setup          func(repository *MockProjectRepository)
		expectedResult *ProjectListResult
		expectedError  error
	}{
		{
			name: "normal case: find projects",
			setup: func(repository *MockProjectRepository) {
				repository.EXPECT().CountProjects(ctx, "user-id").Return(3, nil)
				repository.EXPECT().FindProjects(ctx, "user-id", 1, 1).Return(projects, nil)
			},
			expectedResult: &ProjectListResult{
				Projects: projects,
				PageInfo: PageInfo{
					TotalCount:  3,
					HasPrevious: true,
					HasNext:     true,
				},
			},
			expectedError: nil,
		},
		{
			name: "error case: count projects failed",
			setup: func(repository *MockProjectRepository) {
				repository.EXPECT().CountProjects(ctx, "user-id").Return(0, countProjectsError)
			},
			expectedResult: nil,
			expectedError:  countProjectsError,
		},
		{
			name: "error case: find projects failed",
			setup: func(repository *MockProjectRepository) {
				repository.EXPECT().CountProjects(ctx, "user-id").Return(3, nil)
				repository.EXPECT().FindProjects(ctx, "user-id", 1, 1).Return(nil, findProjectsError)
			},
			expectedResult: nil,
			expectedError:  findProjectsError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repository := NewMockProjectRepository(ctrl)
			test.setup(repository)

			projectUsecase := NewProjectUsecase(repository)

			result, err := projectUsecase.FindProjects(ctx, "user-id", 2, 1)
			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}
