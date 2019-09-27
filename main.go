package main

import (
	"fmt"
	"gopro/app/i18n"
)

func main() {
	fmt.Println(i18n.Trans("tw", "hello.2", map[string]string{
		"Name": "ok",
	}))
}
