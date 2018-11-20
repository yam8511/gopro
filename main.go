package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World")
	r := gin.Default()
	r.Any("/*url", func(c *gin.Context) {
		log.Println("****************************")
		log.Println("        ！機密資料！     ")
		log.Println("****************************")
		err := c.Request.ParseForm()
		if err != nil {
			log.Println("解析 Form 失敗, ", err.Error())
		} else {
			log.Println("form ****", c.Request.Form)
		}
		dd, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("解析 Body 失敗, ", err.Error())
		} else {
			log.Println("body ****", string(dd))
		}
		log.Println("****************************")
		log.Println("        ！機密資料！     ")
		log.Println("****************************")
		URL := c.Param("url")
		log.Println("URL --->", URL)

		referer := c.Request.Header.Get("referer")
		log.Println("Referer --->", referer)
		if referer != "" {
			refURL, err := url.Parse(referer)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			log.Println("Ref Path ===>", refURL.Path)

			refURL, err = url.Parse(normalizeURL(refURL.Path))
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			log.Println("Ref Host ===>", refURL.Scheme, refURL.Host)
			URL = refURL.Scheme + "://" + refURL.Host + URL
		} else {
			URL = normalizeURL(URL)
			log.Println("Normalize URL --->", URL)
		}

		req, err := http.NewRequest(c.Request.Method, URL, nil)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		req.Header.Set("referer", URL)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		b := new(writer)
		defer res.Body.Close()
		b.body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		b.header = res.Header
		c.Render(http.StatusOK, b)
	})

	r.Run(":8000")
}

func normalizeURL(URL string) string {
	log.Println("Raw URL --->", URL)
	if strings.HasPrefix(URL, "/") {
		URL = strings.TrimPrefix(URL, "/")
	}

	if !strings.HasPrefix(URL, "http") {
		URL = "https://" + URL
	}

	return URL
}

type writer struct {
	body   []byte
	header http.Header
}

// Render 渲染
func (c *writer) Render(w http.ResponseWriter) error {

	log.Println("渲染, ", w)
	c.WriteContentType(w)
	_, err := w.Write(c.body)
	return err
}

// WriteContentType 設定回傳型態
func (c *writer) WriteContentType(w http.ResponseWriter) {
	log.Println("設定回傳型態, ", w)
	w.Header().Set("Content-Type", c.header.Get("Content-Type"))
}
