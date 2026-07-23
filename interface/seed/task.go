package seed

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/repository"
	"context"
)

type TaskSeeder struct{}

func NewTaskSeeder() *TaskSeeder {
	return &TaskSeeder{}
}

func (seeder *TaskSeeder) Run(ctx context.Context, client *ent.Client) error {
	taskRepository := repository.CreateTaskRepository(client)

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
