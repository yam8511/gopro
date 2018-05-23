package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"
)

func main() {
	n := time.Now()
	t := n.Format("2006-01-02T15:04:05Z07:00")
	t2 := n.Format(time.RFC3339)
	// hash := []byte{}
	hash := fmt.Sprintf("%x", md5.Sum([]byte(t)))
	hash2 := fmt.Sprintf("%x", md5.Sum([]byte(t2)))
	log.Printf("Time -> %s , Hash -> %s", t, hash)
	log.Printf("Time -> %s , Hash -> %s", t2, hash2)
}
