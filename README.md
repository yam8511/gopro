# Swagger API Demo

0. 下載範例專案
```shell
$ cd $GOPATH/src/
# or
# cd ~/go/src/
$ git clone https://github.com/yam8511/gopro.git -b gin-swagger
$ cd gopro/
```

1. 安裝套件
```shell
$ go get github.com/swaggo/swag/cmd/swag
$ go get -u github.com/swaggo/gin-swagger
$ go get -u github.com/swaggo/gin-swagger/swaggerFiles
```

2. 在 gin.Handler 的 func 上方寫入註解
```go
// @Summary 登入
// @Description 設定cookie
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {string} string "登入成功"
// @Router /login [get]
func login(c *gin.Context) {
    // ...
}
```

3. 產生文檔
```shell
$ swag init
```

4. main 引入文檔
```go
package main
import (
    // ...
	_ "gopro/docs"
)
```

5. 開啟server
see http://127.0.0.1:8000/swagger/index.html

### 參考文章
- [使用swaggo自动生成Restful API文档](https://ieevee.com/tech/2018/04/19/go-swag.html)
- [註解格式說明](https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html)
