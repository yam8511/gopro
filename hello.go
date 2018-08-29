/**
 * ===================================
 *              Go Mod 初探
 * ===================================
 * 版本： go 1.11
 * 功能：
 * 	golang 內建套件管理工具，當進行程式編譯時，會自動下載所需套件。
 * 	即使專案位置是在 GOPATH 之外，也是個可以自動下載所需套件的工具。
 *
 * 環境變數： GO111MODULE
 * 	說明：用來偵測是否啟用 mod 功能
 * 	變數值：
 * 		- on (啟動，編譯時，會自動偵測所需套件並下載)
 * 		- auto (半啟動，會偵測是否有 go.mod，若有會執行 on 的事情)
 * 		- off (不啟動)
 *
 * 參數：
 * 	- init
 * 		初始化，但需要在 GOPATH 底下，
 * 		或是專案底下若有設定 git remote origin，
 * 		會自動抓 origin 的專案名稱當 go.mod 的 module
 * 		ps. 若非在GOPATH底下，無法使用 go mod init，但仍可自己手動建立 go.mod，也有效
 * 	- tidy
 * 		自動增加遺漏的套件，以及移除未使用的套件
 * 	- graph
 * 		在終端機(Terminal)上顯示 go.mod 的內容
 * 	- vendor
 * 		在專案目錄底下，根據 go.mod 產生出 vendor 的資料夾
 * 		以便之後若使用其他套件管理工具，也仍可使用
 * 	- why
 * 		顯示專案需要哪些套件，並自動增加遺漏的套件，但不會移除未使用的套件
 * 	- verify
 * 		檢驗 go.mod 的內容或套件是否有問題
 * 		但使用的體驗上，感覺 tidy、why 也會自動 verify
 * 	- edit
 * 		[待補]
 * 		編輯go.mod文件，選項有-json、-require和-exclude，可以使用幫助go help mod edit
 * 	- download
 * 		[待補]
 * 		下載modules到本地cache
 *
 *
 * 心得：
 * 	其實眾多套件管理工具，個人偏好 govendor，可方便找到檢查專案所需要套件，
 * 	並且可由既有 go get 下載過的套件，再複製進去 vendor，而不會再另外下載，
 * 	以及如有發現遺漏的套件，也會自動額外下載。
 * 	這次初探的情況之下，個人覺得 mod 其實還算蠻方便的，與 govendor 相似。
 * 	可搭配 go mod tidy 與 go mod vendor 做出與原本 govendor 同樣的結果。
 * 	但 mod 更加方便的是，不必在 GOPATH 底下才能使用，而是依賴於 go.mod，
 * 	而且 mod 又是內建的，可不必在額外下載 govendor 工具。
 * 	初探之下，對 mod 頗有好感，但目前仍在實驗階段，期待之後釋出穩定版本。
 *
 * 參考：
 * [Golang modules 初探
](https://tw.saowen.com/a/280205cad1502193482905232ba84e501998fcb7216d8e51b037ca892cb22337)
 *
**/

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
