package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func update(screen *ebiten.Image) error {
	ebitenutil.DebugPrint(screen, "Hello world!")
	// 當「按鍵上」被按下時⋯⋯
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		ebitenutil.DebugPrint(screen, "You're pressing the 'UP' button.")
	}
	// 當「按鍵下」被按下時⋯⋯
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		ebitenutil.DebugPrint(screen, "\nYou're pressing the 'DOWN' button.")
	}
	// 當「按鍵左」被按下時⋯⋯
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		ebitenutil.DebugPrint(screen, "\n\nYou're pressing the 'LEFT' button.")
	}
	// 當「按鍵右」被按下時⋯⋯
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		ebitenutil.DebugPrint(screen, "\n\n\nYou're pressing the 'RIGHT' button.")
	}
	return nil
}

func main() {
	ebiten.Run(update, 320, 240, 2, "Hello world!")
	fmt.Println("OKOKO")
}
