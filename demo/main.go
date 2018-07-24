package main

import (
	"log"
)

func main() {
	log.Println("Hello Demo")
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte(`Hello Drone`))
	// })
	// err := http.ListenAndServe(":8888", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// GetDemo 範利用
func GetDemo() string {
	return "demo"
}
