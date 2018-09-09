package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/skip2/go-qrcode"
	qrdecode "github.com/tuotoo/qrcode"
)

func main() {
	// err := qrcode.WriteFile(
	// 	"https://tofu-rosa.herokuapp.com/graphql",
	// 	qrcode.Medium,
	// 	256,
	// 	"qr.png",
	// )
	// if err != nil {
	// 	fmt.Println("write error")
	// }

	err := qrcode.WriteColorFile(
		"https://tofu-rosa.herokuapp.com/graphql",
		qrcode.Medium,
		256,
		color.White,
		color.RGBA{10, 100, 250, 255},
		"qr.png",
	)
	if err != nil {
		fmt.Println("write error")
	}

	fi, err := os.Open("qr.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fi.Close()
	qrmatrix, err := qrdecode.Decode(fi)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(qrmatrix.Content)
}
