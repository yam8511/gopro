package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World")
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello E.F.K")
	})

	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello API")
	})

	r.Run()
}
