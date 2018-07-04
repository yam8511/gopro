package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorhill/cronexpr"
)

func main() {
	fmt.Println("Hello World")
	now := time.Now()
	parsedTime, parseErr := cronexpr.Parse("*/5 * * * * * *")
	if parseErr != nil {
		log.Fatal("Parse Error ->", parseErr)
		return
	}
	nextTime := parsedTime.NextN(now, 2)
	lastTime := time.Unix(nextTime[0].Unix()-(nextTime[1].Unix()-nextTime[0].Unix()), 0)
	fmt.Println("Now ->", now)
	fmt.Println("Last ->", lastTime)
	fmt.Println(nextTime)
}
