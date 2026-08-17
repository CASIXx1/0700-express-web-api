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

	jwtSecret := os.Getenv("JWT_SECRET")

	userRepository := repository.NewUserRepository(dbClient)
	projectRepository := repository.NewProjectRepository(dbClient)
	tokenService := auth.NewTokenService(jwtSecret)
	authUsecase := usecase.NewAuthUsecase(userRepository, auth.NewPasswordVerifier(), tokenService)
	userUsecase := usecase.NewUserUsecase(userRepository, tokenService)
	projectUsecase := usecase.NewProjectUsecase(projectRepository, tokenService)
	authHandler := handler.NewHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	projectHandler := handler.NewProjectHandler(projectUsecase)

	r.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/signup", authHandler.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/users/me", userHandler.Me).Methods(http.MethodGet)
	r.HandleFunc("/users/projects", projectHandler.FindProjects).Methods(http.MethodGet)
	r.HandleFunc("/users/projects/{slug}", projectHandler.FindProjectBySlug).Methods(http.MethodGet)

	r.Use(middleware.RequestLog)

	s := r.NewRoute().Subrouter()
	s.Use(middleware.Auth)
	s.HandleFunc("/users/me", userHandler.Me).Methods(http.MethodGet)
	s.HandleFunc("/users/projects", projectHandler.FindProjects).Methods(http.MethodGet)
	s.HandleFunc("/users/projects/{slug}", projectHandler.FindProjectBySlug).Methods(http.MethodGet)

	addr := ":8080"
	log.Printf("server listening on %s\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
