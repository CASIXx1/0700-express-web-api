package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"0700-express-web-api/ent"
	"0700-express-web-api/ent/user"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	if err := seedUser(ctx, dbClient); err != nil {
		log.Fatal(err)
	}

	if err := seedProject(ctx, dbClient, "programming"); err != nil {
		log.Fatal(err)
	}

	if err := seedProject(ctx, dbClient, "design"); err != nil {
		log.Fatal(err)
	}

	if err := seedProject(ctx, dbClient, "english"); err != nil {
		log.Fatal(err)
	}
}

func seedUser(ctx context.Context, client *ent.Client) error {
	exists, err := client.User.
		Query().
		Where(user.EmailEQ("admin@example.com")).
		Exist(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return client.User.
		Create().
		SetUsername("admin").
		SetEmail("admin@example.com").
		SetPassword(string(hashedPassword)).
		Exec(ctx)
}

func seedProject(ctx context.Context, client *ent.Client, slug string) error {
	return client.Project.
		Create().
		SetSlug(slug).
		Exec(ctx)
}

//_, err = sqlDB.Exec(`
//	INSERT INTO tasks (title, project_tasks)
//	SELECT task.title, projects.id
//	FROM (
//	  VALUES
//		('Learn Go', 'programming'),
//		('Learn English', 'english'),
//		('Learn Design', 'design')
//	) AS task(title, project_slug)
//	JOIN projects ON projects.slug = task.project_slug;
//`)
//if err != nil {
//	log.Fatal(err)
//}
