package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	defer sqlDB.Close()

	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
		    status VARCHAR(255) NOT NULL DEFAULT 'active'
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlDB.Exec(`
		INSERT INTO users (username, email, password)
		VALUES ('admin', 'admin@example.com', 'password')
		ON CONFLICT (email) DO NOTHING;
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
		    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		    slug VARCHAR(255) NOT NULL UNIQUE
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlDB.Exec(`
		INSERT INTO projects (slug)
		VALUES ('programming'),
		('english'),
		('design')
		ON CONFLICT (slug) DO NOTHING;
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
		    title VARCHAR(255) NOT NULL,
		    status VARCHAR(255) NOT NULL DEFAULT 'scheduled',
		    project_id UUID NOT NULL,
		    FOREIGN KEY (project_id) REFERENCES projects(id),
			UNIQUE (title, project_id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlDB.Exec(`
		INSERT INTO tasks (title, project_id)
		SELECT task.title, projects.id
		FROM (
		  VALUES
			('Learn Go', 'programming'),
			('Learn English', 'english'),
			('Learn Design', 'design')
		) AS task(title, project_slug)
		JOIN projects ON projects.slug = task.project_slug;
	`)
	if err != nil {
		log.Fatal(err)
	}
}
