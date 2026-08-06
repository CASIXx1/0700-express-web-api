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
	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to encode response: %v", err)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(`{"message":"failed to encode response"}`))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_, _ = writer.Write(body)
}
