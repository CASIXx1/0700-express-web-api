package seed

import (
	"context"

	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"
)

type taskSeeder struct{}

func newTaskSeeder() seeder {
	return &taskSeeder{}
}

func (seeder *taskSeeder) Run(ctx context.Context, client *ent.Client) error {
	taskRepository := repository.NewTaskRepository(client)

	if err := taskRepository.CreateTask(ctx, "Learn Go", "programming"); err != nil {
		return err
	}
	if err := taskRepository.CreateTask(ctx, "Learn English", "english"); err != nil {
		return err
	}
	if err := taskRepository.CreateTask(ctx, "Learn Design", "design"); err != nil {
		return err
	}

	return nil
}
