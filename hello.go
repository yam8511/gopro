package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	// 視窗配置
	cfg := pixelgl.WindowConfig{
		Title:     "Gopher Rocks!",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
		VSync:     true,
	}

	// 建立視窗
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		if win.Pressed(pixelgl.MouseButtonLeft) {
			// 關閉視窗
			win.SetClosed(true)
		}

		// Clear 像是重新粉刷牆壁的感覺
		win.Clear(colornames.Skyblue)
		// Update 就是顯示用的
		win.Update()
	}

	// time.Sleep(time.Second * 5)
	// fmt.Println("Close Window")
}

func main() {
	pixelgl.Run(run)
}
