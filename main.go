package main

import (
	"0700-express-web-api/internal/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	// PostgreSQLへの接続
	// ルーティング
	// JSONでのレスポンスの書き込みの確認
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"status": "ng",
				"error":  err.Error(),
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"status": "ng",
				"error":  err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}).Methods(http.MethodGet)

	jwtSecret := os.Getenv("JWT_SECRET")

	authHandler := auth.CreateHandler(db, jwtSecret)
	r.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)

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
