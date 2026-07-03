package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	addr := ":8080"
	log.Printf("server listening on %s\n", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
