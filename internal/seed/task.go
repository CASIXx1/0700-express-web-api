package seed

import (
	"context"
	"time"

	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"
)

type taskSeeder struct{}

func newTaskSeeder() seeder {
	return &taskSeeder{}
}

func (seeder *taskSeeder) Run(ctx context.Context, client *ent.Client) error {
	taskRepository := repository.NewTaskRepository(client)

	jst := time.FixedZone("JST", 9*60*60)
	startingAt := time.Date(2026, 8, 1, 0, 0, 0, 0, jst)
	startedAt := time.Date(2026, 8, 5, 0, 0, 0, 0, jst)
	deadline := time.Date(2026, 9, 30, 0, 0, 0, 0, jst)

	if err := taskRepository.CreateTask(ctx, repository.CreateTaskInput{
		Title:       "Learn Golang",
		Description: "variables, types, functions",
		Status:      "scheduled",
		ProjectSlug: "programming",
		StartingAt:  &startingAt,
		StartedAt:   &startedAt,
		FinishedAt:  nil,
		ArchivedAt:  nil,
		Deadline:    &deadline,
	}); err != nil {
		return err
	}
	if err := taskRepository.CreateTask(ctx, repository.CreateTaskInput{
		Title:       "Learn English",
		Description: "grammar, pronounce, idiom, conversation",
		Status:      "scheduled",
		ProjectSlug: "english",
		StartingAt:  &startingAt,
		StartedAt:   &startedAt,
		FinishedAt:  nil,
		ArchivedAt:  nil,
		Deadline:    &deadline,
	}); err != nil {
		return err
	}
	if err := taskRepository.CreateTask(ctx, repository.CreateTaskInput{
		Title:       "Learn Design",
		Description: "UI, UX",
		Status:      "scheduled",
		ProjectSlug: "design",
		StartingAt:  &startingAt,
		StartedAt:   &startedAt,
		FinishedAt:  nil,
		ArchivedAt:  nil,
		Deadline:    &deadline,
	}); err != nil {
		return err
	}

	return nil
}
