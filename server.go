package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
				<head><title>GopherJS</title></head>
				<body>
					<h1>Hello World</h1>
					<h3>1. go get -v -u  github.com/gopherjs/gopherjs</h1>
					<h3>2. gopherjs build hello.go</h1>
					<h3>3. go run server.go</h1>
					<script src="/hello.js.map" ></script>
					<script src="/hello.js" ></script>
				</body>
			</html>
		`))
	})
	http.HandleFunc("/hello.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./hello.js")
	})
	http.HandleFunc("/hello.js.map", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./hello.js.map")
	})
	http.ListenAndServe(":8000", nil)
}
