package main

import (
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
)

var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1 style="color:red;">Welcome, Ginner!</h1>
</body>
</html>
`))

func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.SetHTMLTemplate(html)

	r.GET("/", func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
			log.Println("推送成功！")
		}
		log.Println("Proto ---> ", c.Request.Proto)
		log.Println("Is TLS ---> ", c.Request.TLS != nil)
		c.HTML(200, "https", gin.H{
			"status": "success",
		})
	})
	go func() {
		err := r.Run(":8001")
		log.Fatal(err)
	}()
	err := r.RunTLS(":8000", "./server.crt", "./server.key")
	// err := autotls.Run(r, "zuolar.me")
	log.Fatal(err)
}
