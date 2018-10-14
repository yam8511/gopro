package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Hello Demo")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := "Hello Demo API at " + time.Now().Format(time.RFC3339)
		log.Println(msg)
		w.Write([]byte(msg))
	})
	http.ListenAndServe(":8888", nil)
}

// GetDemo 範利用
func GetDemo() string {
	return "demo"
}
