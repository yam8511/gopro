package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	site := os.Getenv("SITE")
	if site == "" {
		site = "web"
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		if site == "web" {
			c.JSON(http.StatusOK, "Hello Web")
		} else {
			c.JSON(http.StatusOK, "Hello "+site)
		}
	})

	r.GET("/api/create", func(c *gin.Context) {
		db, err := newDBConnection()
		if err != nil {
			c.JSON(http.StatusOK, "error ---> "+err.Error())
			return
		}
		defer db.Close()

		user := &User{
			Name: "Zuolar",
		}
		err = db.Create(&user).Error
		if err != nil {
			c.JSON(http.StatusOK, "error ---> "+err.Error())
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.GET("/api", func(c *gin.Context) {
		if site == "web" {
			data, err := callAPI("demo")
			if err != nil {
				c.JSON(http.StatusOK, "Web Call API ---> "+err.Error())
				return
			}
			c.JSON(http.StatusOK, "Web Call API ---> "+string(data))
		} else {
			data, err := callAPI("web")
			if err != nil {
				c.JSON(http.StatusOK, "Demo Call API ---> "+err.Error())
				return
			}
			c.JSON(http.StatusOK, "Demo Call API ---> "+string(data))
		}
	})

	port := ":8000"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	fmt.Println("Port --->", port)

	if port != ":8000" {
		go func() {
			r.Run(":8000")
		}()
	}
	r.Run(port)
}

func callAPI(ip string) ([]byte, error) {

	url := fmt.Sprintf("http://%s:8000", ip)
	fmt.Println("URL --->", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
