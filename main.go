package main

import (
	"fmt"
	"log"
	"regexp"
)

func main() {
	reg, err := regexp.Compile("/posts/[0-9]+")
	if err != nil {
		log.Fatal(err)
		return
	}

	ok := reg.MatchString("/posts/123")
	fmt.Println("Hello World ---> ", ok)
}
