package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World")
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.GET("/*any", static.Serve("/", static.LocalFile("./upload/picture", true)))

	r.POST("/", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			fmt.Println("error -> ", err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		files, ok := form.File["file"]
		if ok {
			var size int64
			var name []string
			for i := range files {
				file := files[i]
				size += file.Size
				name = append(name, file.Filename)
				f, err := file.Open()
				if err != nil {
					c.JSON(http.StatusInternalServerError, err.Error())
					return
				}
				defer f.Close()
				fileBytes, err := ioutil.ReadAll(f)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err.Error())
					return
				}
				filename := "./upload/picture/" + file.Filename
				uploadDir := filepath.Dir(filename)
				// fmt.Println("dir --> ", uploadDir)
				// 檢查資料夾是否存在，否則建立資料夾
				_, err = os.Stat(uploadDir)
				if os.IsNotExist(err) {
					err = os.MkdirAll(uploadDir, 0777)
					if err != nil {
						c.JSON(http.StatusInternalServerError, err.Error())
						return
					}
				}
				// 移除原本檔案
				err = os.Remove(filename)
				if os.IsExist(err) {
					if err != nil {
						c.JSON(http.StatusInternalServerError, "rm: "+err.Error())
						return
					}
				}

				// 建立新檔案
				target, err := os.Create(filename)
				if err != nil {
					c.JSON(http.StatusInternalServerError, "create: "+err.Error())
					return
				}
				defer target.Close()

				// 將上傳檔案寫入新檔案
				target.Write([]byte{})
				_, err = target.Write(fileBytes)
				if err != nil {
					c.JSON(http.StatusInternalServerError, "write: "+err.Error())
					return
				}

				// 將檔案權限改成可寫可讀
				err = target.Chmod(0666)
				if err != nil {
					c.JSON(http.StatusInternalServerError, "chmod: "+err.Error())
					return
				}
			}
			c.JSON(http.StatusOK, map[string]interface{}{
				"message": "檔案上傳成功",
				"size":    size,
				"name":    name,
			})
			return
		}
		c.JSON(http.StatusOK, "沒有檔案上傳")
	})
	r.Run(":8000")
}
