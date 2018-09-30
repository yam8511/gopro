package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	for i := 0; i < 255; i++ {
		go func(i int) {
			ip := "10.42.0." + strconv.Itoa(i)
			cmd := exec.Command("ping", ip)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Start()
			time.Sleep(time.Second * 5)
			cmd.Process.Kill()
		}(i)
	}
	time.Sleep(time.Second * 6)
	fmt.Println("Hello World")
}
