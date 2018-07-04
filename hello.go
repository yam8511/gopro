package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	second := flag.Int("s", 1, "Ticker Second")
	flag.Parse()
	ticker := time.NewTicker(time.Second * time.Duration(*second))
	log.Println("Di Da")
	for range ticker.C {
		log.Println("Di Da")
	}
}
