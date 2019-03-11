package main

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Game 遊戲
type Game struct {
	w, h  int
	scale float64
	title string
	bg    *ebiten.Image
}

// Run 啟動遊戲
func (g *Game) Run() error {
	ebiten.SetCursorVisible(true)
	bg, err := ebitenutil.NewImageFromURL("https://upload.wikimedia.org/wikipedia/commons/thumb/2/23/Go_Logo_Aqua.svg/1200px-Go_Logo_Aqua.svg.png")
	if err != nil {
		fmt.Println(err)
		return err
	}
	g.bg = bg
	g.h, g.w = ebiten.ScreenSizeInFullscreen()
	g.scale = 1
	g.title = "Hello World"
	return ebiten.Run(g.render, g.h, g.w, g.scale, g.title)
}

func (g *Game) render(screen *ebiten.Image) error {
	if !ebiten.IsFullscreen() {
		ebiten.SetScreenSize(ebiten.ScreenSizeInFullscreen())
	}

	// 先畫圖
	c := color.NRGBA{0x33, 0x33, 0x33, 0xFF}
	screen.Fill(c)

	opts := &ebiten.DrawImageOptions{}
	// opts.SourceRect.Max = image.Point{320, 240}
	// opts.SourceRect.Min = image.Point{0, 0}
	// opts.GeoM.Translate(float64(x), float64(y))
	w, h := g.bg.Size()
	opts.GeoM.Scale(
		320/float64(w),
		320/float64(h),
	)

	// 渲染 square 畫布到 screen 主畫布上並套用空白選項
	screen.DrawImage(g.bg, opts)

	// 在印字
	FPS := fmt.Sprintf("\n\n\n\n\nFPS: %f", ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, FPS)

	// 從 CursorPosition() 取得 X, Y 座標
	x, y := ebiten.CursorPosition()

	// 顯示一段「X: xx, Y: xx」格式文字
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

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

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("EXIT")
	}

	// 當「滑鼠左鍵」被按下時⋯⋯
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		square, err := ebiten.NewImage(16, 16, ebiten.FilterNearest)
		if err != nil {
			return err
		}
		square.Fill(color.White)

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
	fmt.Println("Hello World")
	game := new(Game)
	err := game.Run()
	fmt.Println(err)
}
