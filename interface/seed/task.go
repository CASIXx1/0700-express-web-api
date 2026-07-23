package seed

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/ent/project"
	"context"
	"log"
)

type TaskSeeder struct{}

func NewTaskSeeder() *TaskSeeder {
	return &TaskSeeder{}
}

func (seeder *TaskSeeder) Run(ctx context.Context, client *ent.Client) error {
	if err := createTask(ctx, client, "Learn Go", "programming"); err != nil {
		log.Fatal(err)
	}
	if err := createTask(ctx, client, "Learn English", "english"); err != nil {
		log.Fatal(err)
	}
	if err := createTask(ctx, client, "Learn Design", "design"); err != nil {
		log.Fatal(err)
	}

	return nil
}

func createTask(ctx context.Context, client *ent.Client, title string, projectSlug string) error {
	return client.Task.
		Create().
		SetTitle(title).
		SetProject(client.Project.
			Query().
			Where(project.SlugEQ(projectSlug)).
			OnlyX(ctx)).
		Exec(ctx)
}
