package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/signintech/gopdf"
)

// 網格數
const grids = 28

// 中英文的位移量
const offsetZhEn = 0.05

//595.28, 841.89 = A4
var rect = &gopdf.Rect{W: 595.28, H: 841.89}

// 網格的X軸長度
var gridX = rect.W / float64(grids)

// 網格的Y軸長度
var gridY = rect.H / float64(grids)

// 字高度(單位：網格)
var fontH = 0.6

// 字寬度(單位：網格)
var fontW = 0.6632

func main() {
	// 開始製作PDF
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *rect})
	defer pdf.Close()

	// 加入字體
	err := pdf.AddTTFFont("zh", "kai.ttf")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = pdf.AddTTFFont("en", "./TTF/Crimson-Roman.ttf")
	if err != nil {
		log.Fatal(err.Error())
	}

	// 印出整頁
	PrintPDF(pdf, "這是一大串標題", "9453", "人仁科技", "2018/11/30")

	// 產生PDF檔案
	err = pdf.WritePdf("hello.pdf")
	log.Println("OK,", err)
}

// PrintPDF 印出PDF
func PrintPDF(
	pdf *gopdf.GoPdf,
	title, siteCode, siteName, auditDate string,
) {
	// 印出整頁
	PrintPageLayout(pdf, title, siteCode, siteName, auditDate)

	count := 13 // 缺失數量
	var y, h float64 = 5.5, 4.5
	for i := 1; i <= count; i++ {
		PrintTableRow(
			pdf, i, y, h,
			"出入口未設置安全通道設施，停工罰款項目？是嗎，應該吧",
			"是嗎，應該吧",
			"snoopy.png", "snoopy.png",
			time.Now(),
		)
		y += h
		if i%5 == 0 && i != count {
			y = 5.5
			// 印出整頁
			PrintPageLayout(pdf, title, siteCode, siteName, auditDate)
		}
	}
}

// PrintTableRow 印出一筆資料
func PrintTableRow(
	pdf *gopdf.GoPdf,
	i int,
	y, h float64,
	punchDescription, improvedDescription,
	punchPhoto, improvedPhoto string,
	improvedTime time.Time,
) {
	// 內文字體 14
	tableBodyCellOpt := gopdf.CellOption{
		Align:  gopdf.Middle | gopdf.Center,                         // Middle and Center
		Border: gopdf.Top | gopdf.Right | gopdf.Bottom | gopdf.Left, // ALL
	}
	pdf.SetLineWidth(1)
	pdf.SetLineType("")

	// 項次
	ChangeFont(pdf, "en", 14)
	ChangeCursorPoint(pdf, 2, y)
	pdf.CellWithOption(GetRect(2, h), strconv.Itoa(i), tableBodyCellOpt)

	// 缺失說明
	ChangeFont(pdf, "zh", 14)
	ChangeCursorPoint(pdf, 4, y)
	PrintTableText(pdf, 4, h, 4, y, punchDescription, &tableBodyCellOpt)

	// 缺失照片
	ChangeCursorPoint(pdf, 8, y)
	punchPhoto = strings.TrimSpace(punchPhoto)
	if punchPhoto != "" {
		pdf.Image(punchPhoto, GetPointX(8), GetPointY(y), GetRect(6, h))
		pdf.CellWithOption(GetRect(6, h), "", tableBodyCellOpt)
	} else {
		ChangeFont(pdf, "zh", 14)
		pdf.CellWithOption(GetRect(6, h), "未提供照片", tableBodyCellOpt)
	}

	// 改善說明
	ChangeCursorPoint(pdf, 14, y)
	PrintTableText(pdf, 4, h, 14, y, improvedDescription, &tableBodyCellOpt)

	// 改善照片
	ChangeCursorPoint(pdf, 18, y)
	improvedPhoto = strings.TrimSpace(improvedPhoto)
	if improvedPhoto != "" {
		pdf.Image(improvedPhoto, GetPointX(18), GetPointY(y), GetRect(6, h))
		pdf.CellWithOption(GetRect(6, h), "", tableBodyCellOpt)
	} else {
		ChangeFont(pdf, "zh", 14)
		pdf.CellWithOption(GetRect(6, h), "未提供照片", tableBodyCellOpt)
	}

	// 改善時間
	ChangeCursorPoint(pdf, 24, y)
	ChangeFont(pdf, "en", 14)
	PrintTableTime(pdf, 4, h, 24, y, improvedTime, &tableBodyCellOpt)
}

// PrintTableText 印出表格文字
func PrintTableText(pdf *gopdf.GoPdf, w, h, x, y float64, text string, opt *gopdf.CellOption) {
	// 內文字體 14
	tableBodyCellOpt := gopdf.CellOption{}
	if opt != nil {
		tableBodyCellOpt = *opt
	}

	// 中文是 3 bytes, 英文數字是 1 byte
	txtLen := len(text) / 3
	fontNum := int((w*gridX - 0.4) / (fontW * gridX)) // 一行字數
	lineNum := txtLen / fontNum                       // 有幾行
	if txtLen%fontNum != 0 {
		lineNum++
	}

	if lineNum > 1 {
		pdf.CellWithOption(GetRect(w, h), "", tableBodyCellOpt)
		paddingTop := (4.5 - fontH*float64(lineNum)) / 2 // 距離上方
		// log.Println("字串長度--->", txtLen)
		// log.Println("字數--->", fontNum)
		// log.Println("行數--->", lineNum)
		// log.Println("距離上方--->", paddingTop)
		for i := 0; i < lineNum; i++ {
			txt := ""
			if i == lineNum-1 {
				H := i * 3 * fontNum
				txt = text[H:]
			} else {
				H := i * 3 * fontNum
				T := (i + 1) * 3 * fontNum
				txt = text[H:T]
			}
			ChangeCursorPoint(pdf, x+0.2, y+paddingTop+float64(i)*fontH)
			pdf.Cell(nil, txt)
		}
	} else {
		pdf.CellWithOption(GetRect(w, h), text, tableBodyCellOpt)
	}

}

// PrintTableTime 印出表格文字
func PrintTableTime(pdf *gopdf.GoPdf, w, h, x, y float64, t time.Time, opt *gopdf.CellOption) {
	// 內文字體 14
	tableBodyCellOpt := gopdf.CellOption{
		Align: gopdf.Bottom | gopdf.Center,
	}
	if opt != nil {
		tableBodyCellOpt = *opt
	}

	pdf.CellWithOption(GetRect(w, h), "", tableBodyCellOpt)
	// paddingTop := (4.5 - fontH*2) / 2 // 距離上方
	// log.Println("字串長度--->", txtLen)
	// log.Println("字數--->", fontNum)
	// log.Println("行數--->", lineNum)
	// log.Println("距離上方--->", paddingTop)

	ChangeFont(pdf, "en", 14)
	ChangeCursorPoint(pdf, 24, y)
	pdf.CellWithOption(GetRect(w, h/2-0.1), t.Format("2006/01/02"), gopdf.CellOption{
		Align: gopdf.Bottom | gopdf.Center,
		// Border: gopdf.Bottom,
	})

	ChangeCursorPoint(pdf, 24, y+h/2)
	pdf.CellWithOption(GetRect(w, h/2+0.1), t.Format("15:04:05"), gopdf.CellOption{
		Align: gopdf.Top | gopdf.Center,
		// Border: gopdf.Top,
	})
}

// PrintPageLayout 印出整頁外框
func PrintPageLayout(pdf *gopdf.GoPdf, title, siteCode, siteName, auditDate string) {
	// 新增第一頁
	pdf.AddPage()

	// 先畫網格
	PrintGrids(pdf, rect)

	// 寫入標題,30
	PrintTitle(pdf, title)

	// 寫入工地與稽核,16
	PrintSiteAudit(pdf, siteCode, siteName, auditDate)

	// 畫標頭表格,16
	PrintTableHeader(pdf)
}

// PrintTableHeader 印出表單標頭
func PrintTableHeader(pdf *gopdf.GoPdf) {
	tableHeaderCellOpt := gopdf.CellOption{
		Align:  gopdf.Middle | gopdf.Center,                         // Middle and Center
		Border: gopdf.Top | gopdf.Right | gopdf.Bottom | gopdf.Left, // ALL
	}
	pdf.SetLineType("")
	pdf.SetLineWidth(1.2)
	ChangeFont(pdf, "zh", 16)
	bold := 3
	pointY := 4.5
	for i := 0; i < bold; i++ {
		ChangeCursorPoint(pdf, 2, pointY)
		pdf.CellWithOption(GetRect(2, 1), "項次", tableHeaderCellOpt)
		ChangeCursorPoint(pdf, 4, pointY)
		pdf.CellWithOption(GetRect(4, 1), "缺失說明", tableHeaderCellOpt)
		ChangeCursorPoint(pdf, 8, pointY)
		pdf.CellWithOption(GetRect(6, 1), "缺失照片", tableHeaderCellOpt)
		ChangeCursorPoint(pdf, 14, pointY)
		pdf.CellWithOption(GetRect(4, 1), "改善說明", tableHeaderCellOpt)
		ChangeCursorPoint(pdf, 18, pointY)
		pdf.CellWithOption(GetRect(6, 1), "改善照片", tableHeaderCellOpt)
		ChangeCursorPoint(pdf, 24, pointY)
		pdf.CellWithOption(GetRect(4, 1), "改善時間", tableHeaderCellOpt)
	}
}

// PrintSiteAudit 顯示工地與稽核
func PrintSiteAudit(pdf *gopdf.GoPdf, siteCode, siteName, auditDate string) {
	const siteY = 3.4
	const auditY = siteY + 0.7
	// 寫入工地
	ChangeFont(pdf, "zh", 16)
	ChangeCursorPoint(pdf, 3, siteY)
	pdf.Cell(nil, "工地：")

	ChangeFont(pdf, "en", 16)
	ChangeCursorPoint(pdf, 5.1, siteY+offsetZhEn)
	pdf.Cell(nil, siteCode)

	ChangeFont(pdf, "zh", 16)
	ChangeCursorPoint(pdf, 6.7, siteY)
	pdf.Cell(nil, siteName)

	// 分隔線
	ChangeFont(pdf, "en", 16)
	for i := 0; i < 3; i++ {
		ChangeCursorPoint(pdf, 9.5, siteY+offsetZhEn)
		pdf.Cell(nil, "  |  ")
	}

	// 寫入稽核日期,16
	ChangeFont(pdf, "zh", 16)
	ChangeCursorPoint(pdf, 10.5, siteY)
	pdf.Cell(nil, "稽核日期：")

	ChangeFont(pdf, "en", 16)
	ChangeCursorPoint(pdf, 14, siteY+offsetZhEn)
	pdf.Cell(nil, auditDate)
}

// PrintTitle 印出標題
func PrintTitle(pdf *gopdf.GoPdf, title string) {
	// 寫入標題,30
	ChangeFont(pdf, "zh", 30)
	ChangeCursorPoint(pdf, 0, 1.6)
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
	ChangeFont(pdf, "en", 20)
	ChangeCursorPoint(pdf, 1, 1)
	pdf.SetFillColor(200, 0, 0)
	pdf.CellWithOption(
		&gopdf.Rect{W: 10 * gridX, H: gridY},
		fmt.Sprintf("X: %f, Y:%f", gridX, gridY),
		gopdf.CellOption{
			Align: gopdf.Left,
		},
	)
	pdf.SetLineWidth(0.1)
	pdf.SetLineType("dotted")
	count := grids
	for i := 1; i <= count; i++ {
		for j := 1; j <= count; j++ {
			x1 := float64(i-1) * gridX
			y1 := float64(j-1) * gridY
			x2 := float64(i) * gridX
			y2 := float64(j) * gridY
			pdf.Line(x1, y1, x1, y2)
			pdf.Line(x1, y1, x2, y1)
			pdf.Line(x2, y2, x2, y1)
			pdf.Line(x2, y2, x1, y2)
			// ChangeCursorPoint(pdf, float64(i), float64(j))
			// pdf.CellWithOption(GetRect(1, 1), "", gopdf.CellOption{Border: 63})
		}
	}
	pdf.SetFillColor(0, 0, 0)
}

// GetRect 取矩型
func GetRect(w, h float64) *gopdf.Rect {
	return &gopdf.Rect{
		W: gridX * w,
		H: gridY * h,
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
