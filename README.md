# Go WebAssembly

1. 編譯 wasm 檔案
```shell
$ GOOS=js GOARCH=wasm go build -o test.wasm main.go
```
ps. 其中 html, js 檔案來自 /usr/local/go/
(go env GOROOT)/misc/wasm/wasm_exec.{html,js}

2. 開啟 server 執行看看
```shell
$ go run server.go
```

3. 開啟瀏覽器 -> http://localhost:8080/wasm_exec.html
