package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	s, _ := base64.URLEncoding.DecodeString("SABlAGwAbABvAA==")
	fmt.Println(s, string(s), string(s) == "Hello")
	s = bytes.Replace(s, []byte{0}, nil, -1)
	fmt.Println(s, string(s), string(s) == "Hello")
}
