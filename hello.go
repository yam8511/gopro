package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World")

	// 建立 Router
	r := gin.Default()

	staticFiles := []string{
		"./public/404.html",
	}

	// 定義 PATH
	switch os.Getenv("SITE") {
	case "game":
		r.GET("/css/*file", fileHandler)
		r.GET("/js/*file", fileHandler)
		r.GET("/image/*file", fileHandler)
		r.GET("/game/:id", gameHandler)
	case "admin":
		r.Use(static.Serve("/", static.LocalFile("./public/admin", false)))
		staticFiles = append(staticFiles, "./public/admin/index.html")
		r.GET("/", adminHandler)
	}

	// 載入檔案
	r.LoadHTMLFiles(staticFiles...)
	// 設置404回傳
	r.NoRoute(notFoundHandler)

	// 啟動伺服器
	err := r.Run(":8000")
	log.Fatal(err)
}

func fileHandler(c *gin.Context) {
	refer := c.Request.Referer()
	log.Println("Referer --->", refer)

	// 先預處理網址
	refer = strings.Replace(refer, "//", "/", -1)
	paths := strings.Split(refer, "/")
	log.Println("分解 --->", paths, len(paths))
	if len(paths) < 4 {
		c.Status(404)
		return
	}

	fs := map[string]gin.HandlerFunc{
		"0": static.Serve("/", static.LocalFile("./public/game", false)),
		"1": static.Serve("/", static.LocalFile("./public/game1", false)),
	}
	id := paths[3]
	next, ok := fs[id]
	if !ok {
		c.Status(404)
		return
	}
	next(c)
}

func gameHandler(c *gin.Context) {
	id := c.Param("id")
	whiteList := map[string]string{
		"0": "game",
		"1": "game1",
	}

	if id == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/game/0")
		return
	}

	for gid, path := range whiteList {
		if id == gid {
			if pusher := c.Writer.Pusher(); pusher != nil {
				// use pusher.Push() to do server push
				if err := pusher.Push("./public/"+path+"/css/main.css", nil); err != nil {
					log.Printf("Failed to push: %v", err)
				}
			} else {
				log.Println("不支援 Server Push")
			}

			c.Header("Content-Type", "text/html; charset=utf-8")
			c.File("./public/" + path + "/index.html")
			return
		}
	}
	c.Status(http.StatusNotFound)
}

func adminHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func notFoundHandler(c *gin.Context) {
	c.Status(http.StatusNotFound)
	c.HTML(http.StatusNotFound, "404.html", nil)
}
