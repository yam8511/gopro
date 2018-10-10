package main

import (
	"fmt"
	"time"
)

func main() {
	defer func() {
		if except := recover(); except != nil {
			fmt.Println("主程序接收到panic --->", except)
		}
	}()
	fmt.Println("Hello World")
	go func() {
		defer func() {
			if except := recover(); except != nil {
				fmt.Println("副程序接收到panic --->", except)
			}
		}()
		h()
	}()
	// go h() // recover 無法接收不同 goroutine 的 panic
	time.Sleep(time.Second)
}

func h() {
	panic("Shutdown")
}
