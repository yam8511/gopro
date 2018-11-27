package main

import (
	"fmt"
	"log"

	"github.com/signintech/gopdf"
)

// 網格數
const grids = 28

//595.28, 841.89 = A4
var rect = &gopdf.Rect{W: 595.28, H: 841.89}

// 網格的X軸長度
var gridX = rect.W / float64(grids)

// 網格的Y軸長度
var gridY = rect.H / float64(grids)

func main() {
	// 開始製作PDF
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *rect})

	// 加入字體
	err := pdf.AddTTFFont("zh", "kai.ttf")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = pdf.AddTTFFont("en", "./TTF/Crimson-Roman.ttf")
	if err != nil {
		log.Fatal(err.Error())
	}

	// 新增第一頁
	pdf.AddPage()

	// 寫入標題,30
	PrintTitle(pdf, "這是一大串標題")

	// 寫入工地與稽核,16
	PrintSiteAudit(pdf, "9453", "人仁科技", "2018/11/30")

	// 項次
	// 缺失說明
	// 缺失照片
	// 改善說明
	// 改善照片
	// 改善時間
	// pdf.SetLineWidth(10)
	// pdf.Line(GetPointX(3), GetPointY(6), GetPointX(3), GetPointY(10.5))

	// 先畫網格
	ChangeFont(pdf, "en", 20)
	PrintGrids(pdf, rect)

	// 產生PDF檔案
	err = pdf.WritePdf("hello.pdf")
	log.Println("OK,", err)

	// 照片5張一頁, 有23格去分攤, 1張照片，只能佔 4.6格
	// 表格 16
	// 內文字體 14
}

// PrintSiteAudit 顯示工地與稽核
func PrintSiteAudit(pdf *gopdf.GoPdf, siteCode, siteName, auditDate string) {
	// 寫入工地
	ChangeFont(pdf, "zh", 16)
	ChangeCursorPoint(pdf, 3, 3.8)
	pdf.Cell(nil, "工地：")

	ChangeFont(pdf, "en", 16)
	ChangeCursorPoint(pdf, 5.1, 3.85)
	pdf.Cell(nil, siteCode)

	ChangeFont(pdf, "zh", 16)
	ChangeCursorPoint(pdf, 6.7, 3.8)
	pdf.Cell(nil, siteName)

	// 寫入稽核日期,16
	ChangeFont(pdf, "zh", 16)
	ChangeCursorPoint(pdf, 3, 4.5)
	pdf.Cell(nil, "稽核日期：")

	ChangeFont(pdf, "en", 16)
	ChangeCursorPoint(pdf, 6.6, 4.55)
	pdf.Cell(nil, auditDate)
}

// PrintTitle 印出標題
func PrintTitle(pdf *gopdf.GoPdf, title string) {
	// 寫入標題,30
	ChangeFont(pdf, "zh", 30)
	ChangeCursorPoint(pdf, 0, 2)
	pdf.CellWithOption(
		&gopdf.Rect{W: rect.W, H: GetPointY(2)},
		title,
		gopdf.CellOption{
			Align: gopdf.Center,
		},
	)
}

// PrintGrids 印出網格
func PrintGrids(pdf *gopdf.GoPdf, rect *gopdf.Rect) {
	ChangeCursorPoint(pdf, 1, 1)
	pdf.CellWithOption(
		&gopdf.Rect{W: 10 * gridX, H: gridY},
		fmt.Sprintf("X: %f, Y:%f", gridX, gridY),
		gopdf.CellOption{
			Align: gopdf.Left,
		},
	)
	pdf.SetLineType("dotted")
	count := grids
	for i := 0; i <= count; i++ {
		for j := 0; j < count; j++ {
			x1 := float64(i) * gridX
			y1 := float64(j) * gridY
			x2 := float64((i + 1)) * gridX
			y2 := float64((j + 1)) * gridY
			pdf.Line(x1, y1, x1, y2)
			pdf.Line(x1, y1, x2, y1)
			pdf.Line(x2, y2, x2, y1)
			pdf.Line(x2, y2, x1, y2)
		}
	}
}

// ChangeFont 變更字體
func ChangeFont(pdf *gopdf.GoPdf, font string, size int) {
	err := pdf.SetFont(font, "", size)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

// ChangeCursorPoint 變更游標位置
func ChangeCursorPoint(pdf *gopdf.GoPdf, x, y float64) {
	pdf.SetX(GetPointX(x))
	pdf.SetY(GetPointY(y))
}

// GetPointX 取X位子
func GetPointX(i float64) float64 {
	if i < 1 {
		i = 1
	}
	return gridX * float64(i-1)
}

// GetPointY 取Y位子
func GetPointY(i float64) float64 {
	if i < 1 {
		i = 1
	}
	return gridY * float64(i-1)
}
