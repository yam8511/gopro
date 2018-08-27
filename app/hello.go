package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello World")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	port := ":8000"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	r.Run(port)
}
