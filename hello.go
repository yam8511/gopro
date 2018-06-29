package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	fmt.Println("Hello, playground")
	js.Global.Call("alert", "Hello, JavaScript")
	println("Hello, JS console")
}
