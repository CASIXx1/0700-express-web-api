package seed

import (
	"context"
	"time"

	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"
)

type projectSeeder struct{}

func newProjectSeeder() seeder {
	return &projectSeeder{}
}

func (seeder *projectSeeder) Run(ctx context.Context, client *ent.Client) error {
	projectRepository := repository.NewProjectRepository(client)
	userRepository := repository.NewUserRepository(client)

	user, err := userRepository.FindUserByEmail(ctx, "test@example.com")
	if err != nil {
		return err
	}

	jst := time.FixedZone("JST", 9*60*60)
	deadline := time.Date(2026, 9, 30, 0, 0, 0, 0, jst)
	startingAt := time.Date(2026, 8, 1, 0, 0, 0, 0, jst)

	if err := projectRepository.CreateProject(ctx, repository.CreateProjectInput{
		Name:       "Programming",
		Slug:       "programming",
		Goal:       "Build programming skills",
		Shouldbe:   "It should improve programming ability",
		Color:      "#FF0000",
		Deadline:   deadline,
		StartingAt: startingAt,
		StartedAt:  startingAt,
		FinishedAt: deadline,
		UserID:     user.ID,
		SortOrder:  1,
	}); err != nil {
		return err
	}

	if err := projectRepository.CreateProject(ctx, repository.CreateProjectInput{
		Name:       "English",
		Slug:       "english",
		Goal:       "Improve English communication",
		Shouldbe:   "It should improve English fluency",
		Color:      "#00AAFF",
		Deadline:   deadline,
		StartingAt: startingAt,
		StartedAt:  startingAt,
		FinishedAt: deadline,
		UserID:     user.ID,
		SortOrder:  2,
	}); err != nil {
		return err
	}

	if err := projectRepository.CreateProject(ctx, repository.CreateProjectInput{
		Name:       "Design",
		Slug:       "design",
		Goal:       "Improve product design",
		Shouldbe:   "It should improve design quality",
		Color:      "#00CC66",
		Deadline:   deadline,
		StartingAt: startingAt,
		StartedAt:  startingAt,
		FinishedAt: deadline,
		UserID:     user.ID,
		SortOrder:  3,
	}); err != nil {
		return err
	}

	return nil
}
