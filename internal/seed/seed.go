package seed

import (
	"context"

	"0700-express-web-api/ent"
)

type seeder interface {
	Run(ctx context.Context, client *ent.Client) error
}

type defaultSeeder struct {
	seeders []seeder
}

func NewSeeder() *defaultSeeder {
	return &defaultSeeder{
		seeders: []seeder{
			newUserSeeder(),
			newProjectSeeder(),
			newTaskSeeder(),
		},
	}
}

func (s *defaultSeeder) Run(ctx context.Context, client *ent.Client) error {
	for _, seeder := range s.seeders {
		if err := seeder.Run(ctx, client); err != nil {
			return err
		}
	}

	return nil
}
