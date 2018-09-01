// 參考網址： https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin

package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteString(b string) (int, error) {
	w.body.Write([]byte(b))
	return w.ResponseWriter.Write([]byte(b))
}

func ginBodyLogMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	go func() {
		if c.Writer.Status() == 200 && c.Writer.Written() {
			time.Sleep(time.Second * 3)
			log.Println("Response Body:", blw.body.String())
			// var data map[string]interface{}
			// err := json.Unmarshal(blw.body.Bytes(), &data)
			// if err != nil {
			// 	log.Println("Get Response Data Error ---> " + err.Error())
			// 	return
			// }

			// for key, val := range data {
			// 	log.Printf("Response [%s]: %v, %T \n", key, val, val)
			// 	switch val.(type) {
			// 	case string:
			// 	case bool:
			// 	case int:
			// 	case map[string]interface{}:
			// 		mapVal := val.(map[string]interface{})
			// 		for key2, val2 := range mapVal {
			// 			log.Printf("Response [%s]: %v, %T \n", key2, val2, val2)
			// 		}
			// 	}
			// }
		}
	}()
}

type data struct {
	Name string `json:"name"`
	Age  int
}

func main() {
	r := gin.Default()
	r.Use(ginBodyLogMiddleware)
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status": true,
			"object": data{"Zuolar", 24},
			"result": map[string]string{
				"a": "Apple",
				"z": "Zoo",
			},
		})
	})
	r.GET("/str", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	err := r.Run()
	log.Fatal(err)
}
