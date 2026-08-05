package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type normalResponse[T any] struct {
	Data T `json:"data"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func writeResponse[T any](writer http.ResponseWriter, statusCode int, response T) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
