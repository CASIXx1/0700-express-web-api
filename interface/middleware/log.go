package middleware

import (
	"log"
	"net/http"
)

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("request started: %s %s", request.Method, request.URL.RequestURI())
		next.ServeHTTP(writer, request)
		log.Printf("request finished: %s %s", request.Method, request.URL.RequestURI())
	})
}
