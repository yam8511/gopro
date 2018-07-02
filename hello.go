package main

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func update(screen *ebiten.Image) error {
	// 先畫圖
	c := color.NRGBA{0xff, 0x00, 0x00, 0xff}
	screen.Fill(c)

	// 在印字
	FPS := fmt.Sprintf("\n\n\n\n\nFPS: %f", ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, FPS)

	// 從 CursorPosition() 取得 X, Y 座標
	x, y := ebiten.CursorPosition()

	// 顯示一段「X: xx, Y: xx」格式文字
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

	square, err := ebiten.NewImage(16, 16, ebiten.FilterNearest)
	if err != nil {
		return err
	}
	square.Fill(color.White)

	// opts.GeoM.Rotate(32)

	// 渲染 square 畫布到 screen 主畫布上並套用空白選項
	// screen.DrawImage(square, opts)

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

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("EXIT")
	}

	// 當「滑鼠左鍵」被按下時⋯⋯
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// 從 CursorPosition() 取得 X, Y 座標
		x, y = ebiten.CursorPosition()
		// 建立一個空白選項建構體
		opts := &ebiten.DrawImageOptions{}
		// opts.GeoM.Translate(0, 0)
		opts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(square, opts)
		ebitenutil.DebugPrint(screen, "\nYou're pressing the 'LEFT' mouse button.")
	}
	// 當「滑鼠右鍵」被按下時⋯⋯
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		ebitenutil.DebugPrint(screen, "\nYou're pressing the 'RIGHT' mouse button.")
	}
	// 當「滑鼠中鍵」被按下時⋯⋯
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		ebitenutil.DebugPrint(screen, "\n\nYou're pressing the 'MIDDLE' mouse button.")
	}

	return nil
}

func main() {
	err := ebiten.Run(update, 320, 240, 2, "Hello world!")
	fmt.Println(err)
}
