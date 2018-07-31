package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World")

	r := gin.Default()

	api := r.Group("/api")
	{
		demo := api.Group("demo")
		demo.GET("id", func(c *gin.Context) {
			fmt.Println("Keys ---> ", c.Keys)
			c.Set("name", time.Now().Unix())
			a := c.Query("a")
			fmt.Println("Get Query a --->", a)
			c.Next()
			cd := c.Query("c")
			fmt.Println("Get Query cd --->", cd)
		}, func(c *gin.Context) {
			name := c.GetInt64("name")
			b := c.Query("b")
			fmt.Println("Get Query b --->", b)
			c.JSON(http.StatusOK, name)
			return
		})
	}

	err := r.Run(":8000")
	log.Fatal(err)
}
