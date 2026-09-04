package main

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/handler"
	"0700-express-web-api/interface/middleware"
	"0700-express-web-api/interface/repository"
	"0700-express-web-api/internal/auth"
	"0700-express-web-api/usecase"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

	dbClient, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	// Migration
	if err := dbClient.Schema.Create(context.Background()); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(middleware.RequestLog)

	jwtSecret := os.Getenv("JWT_SECRET")

	userRepository := repository.NewUserRepository(dbClient)
	projectRepository := repository.NewProjectRepository(dbClient)
	taskRepository := repository.NewTaskRepository(dbClient)
	tokenService := auth.NewTokenService(jwtSecret)
	authUsecase := usecase.NewAuthUsecase(userRepository, auth.NewPasswordVerifier(), tokenService)
	userUsecase := usecase.NewUserUsecase(userRepository)
	projectUsecase := usecase.NewProjectUsecase(projectRepository)
	taskUsecase := usecase.NewTaskUsecase(taskRepository)
	authHandler := handler.NewHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	projectHandler := handler.NewProjectHandler(projectUsecase)
	taskHandler := handler.NewTaskHandler(taskUsecase)

	r.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/signup", authHandler.SignUp).Methods(http.MethodPost)

	authRequired := r.NewRoute().Subrouter()
	authRequired.Use(middleware.Auth(tokenService, userRepository))
	authRequired.HandleFunc("/users/me", userHandler.Me).Methods(http.MethodGet)
	authRequired.HandleFunc("/users/projects", projectHandler.FindProjects).Methods(http.MethodGet)
	authRequired.HandleFunc("/users/projects/{slug}", projectHandler.FindProjectBySlug).Methods(http.MethodGet)
	authRequired.HandleFunc("/users/tasks", taskHandler.FindTasks).Methods(http.MethodGet)
	authRequired.HandleFunc("/users/tasks", taskHandler.CreateTask).Methods(http.MethodPost)
	authRequired.HandleFunc("/users/tasks/{id}", taskHandler.FindTaskByID).Methods(http.MethodGet)
	authRequired.HandleFunc("/users/tasks/{id}", taskHandler.DeleteTask).Methods(http.MethodDelete)

	addr := ":8080"
	log.Printf("server listening on %s\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
