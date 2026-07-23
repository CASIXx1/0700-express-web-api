package seed

import (
	"0700-express-web-api/ent"
	"context"
	"log"
)

type ProjectSeeder struct{}

func NewProjectSeeder() *ProjectSeeder {
	return &ProjectSeeder{}
}

func (seeder *ProjectSeeder) Run(ctx context.Context, client *ent.Client) error {
	if err := CreateProject(ctx, client, "programming"); err != nil {
		log.Fatal(err)
	}

	if err := CreateProject(ctx, client, "design"); err != nil {
		log.Fatal(err)
	}

	if err := CreateProject(ctx, client, "english"); err != nil {
		log.Fatal(err)
	}

	return nil
}

func CreateProject(ctx context.Context, client *ent.Client, slug string) error {
	return client.Project.
		Create().
		SetSlug(slug).
		Exec(ctx)
}
