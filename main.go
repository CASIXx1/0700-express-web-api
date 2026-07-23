package main

import (
	"0700-express-web-api/ent"
	"0700-express-web-api/interface/handler"
	"0700-express-web-api/interface/repository"
	"0700-express-web-api/usecase"
	"context"
	"encoding/json"
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

	// PostgreSQLへの接続
	// ルーティング
	// JSONでのレスポンスの書き込みの確認
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}).Methods(http.MethodGet)

	jwtSecret := os.Getenv("JWT_SECRET")

	authRepository := repository.CreateAuthRepository(dbClient)
	authUsecase := usecase.CreateAuthUsecase(authRepository, jwtSecret)
	authHandler := handler.CreateHandler(authUsecase)

	r.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/signup", authHandler.SignUp).Methods(http.MethodPost)

	addr := ":8080"
	log.Printf("server listening on %s\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
