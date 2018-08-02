package main

import (
	"fmt"
	"log"

	_ "gopro/docs"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	// _ "github.com/swaggo/swag/example/celler/docs"
)

// @title GoPro API
// @version 0.0.1
// @description  範例API文件工具.
// @termsOfService http://swagger.io/terms/
// @license.name Zuolar
func main() {
	fmt.Println("Hello World")

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.POST("/login", login)

		api.GET("/user/:id", getUser)

		api.PUT("/password", updatePasswd)

		api.DELETE("/logout", logout)
	}

	err := r.Run(":8000")
	log.Fatal(err)
}
