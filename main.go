package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello~ " + os.Getenv("APP")))
	})
	http.ListenAndServe(":8000", nil)
}
