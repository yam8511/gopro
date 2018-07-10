package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello World")

	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		return
	}
}
