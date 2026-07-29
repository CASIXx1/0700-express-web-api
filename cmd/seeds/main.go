package main

import (
	"0700-express-web-api/internal/seed"
	"context"
	"fmt"
	"log"
	"os"

	"0700-express-web-api/ent"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	dbClient, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	tx, err := dbClient.Tx(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	seeds := seed.NewSeeder()

	if err := seeds.Run(ctx, tx.Client()); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
