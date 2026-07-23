package seed

import (
	"context"

	"0700-express-web-api/ent"
)

type Seeder interface {
	Run(ctx context.Context, client *ent.Client) error
}

type DefaultSeeder struct {
	seeders []Seeder
}

func NewSeeder(seeders ...Seeder) *DefaultSeeder {
	return &DefaultSeeder{
		seeders: seeders,
	}
}

func (s *DefaultSeeder) Run(ctx context.Context, client *ent.Client) error {
	for _, seeder := range s.seeders {
		if err := seeder.Run(ctx, client); err != nil {
			return err
		}
	}

	return nil
}
