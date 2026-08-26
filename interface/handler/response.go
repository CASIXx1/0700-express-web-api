package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type normalResponse[T any] struct {
	Data T `json:"data"`
}

type paginatedResponse[T any] struct {
	Data     T                  `json:"data"`
	PageInfo paginationResponse `json:"pageInfo"`
}

type paginationResponse struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
	TotalCount  int  `json:"totalCount"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteResponse[T any](writer http.ResponseWriter, statusCode int, response T) {
	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to encode response: %v", err)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)

		if _, err := writer.Write([]byte(`{"message":"failed to encode response"}`)); err != nil {
			panic("failed to write error response: " + err.Error())
		}

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	if _, err := writer.Write(body); err != nil {
		panic("failed to write response: " + err.Error())
	}
}
