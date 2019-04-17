package main

import (
	"net/http/pprof"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.GET("/", index)
	p := r.Group("/debug/pprof")
	p.GET("/", pprofIndex)
	p.GET("/:name", pprofIndex)
	r.Run(":8000")
}

func pprofIndex(c *gin.Context) {
	switch c.Param("name") {
	case "profile":
		pprof.Profile(c.Writer, c.Request)
	default:
		pprof.Index(c.Writer, c.Request)
	}
}

func index(c *gin.Context) {
	c.JSON(200, getData())
}

func getData() []string {
	obj := []string{}
	for i := 0; i < 10000; i++ {
		obj = append(obj, strconv.Itoa(i))
	}
	return obj
}
