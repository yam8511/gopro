# Go WebAssembly

## 編譯 wasm 檔案

```shell
GOOS=js GOARCH=wasm go build -o test.wasm main.go
```

## 複製HTML檔案

html, js 檔案來自 /usr/local/go/

```shell
cp $(go env GOROOT)/misc/wasm/wasm_exec.{html,js} .
```

## 開啟 server 執行看看

```shell
go run ./server/server.go
```

## [開啟瀏覽器](http://localhost:8080/wasm_exec.html)
