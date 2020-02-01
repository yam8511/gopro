package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
)

const (
	dx = 57
	dy = 72
)

func main() {
	CreateRawCardImage()
	CreateCoverPNG()
	CreateCardPNG()
	CreateAllPNG()
	CreateAllGoFile()
}

// CreateAllGoFile 產生所有麻將圖檔的Go檔案
func CreateAllGoFile() {
	allImg := AllImamge()
	for i, imgs := range allImg {
		for j, img := range imgs {
			prefix := ""
			switch i {
			case 1:
				prefix = "B"
			case 2:
				prefix = "L"
			case 3:
				prefix = "O"
			case 4:
				prefix = "T"
			}
			name := prefix + strconv.Itoa(j)
			CreateGoFile("majo/"+name+".go", name+"_PNG", img)
		}
	}
}

// CreateAllPNG 產生所有麻將的圖檔
func CreateAllPNG() {
	allImg := AllImamge()
	for i, imgs := range allImg {
		for j, img := range imgs {
			switch i {
			case 1:
				CreatePNG("majo/B"+strconv.Itoa(j)+".png", img)
			case 2:
				CreatePNG("majo/L"+strconv.Itoa(j)+".png", img)
			case 3:
				CreatePNG("majo/O"+strconv.Itoa(j)+".png", img)
			case 4:
				CreatePNG("majo/T"+strconv.Itoa(j)+".png", img)
			}
		}
	}
}

// AllImamge 所有麻將圖檔
func AllImamge() map[int]map[int]image.Image {
	newCard := CardImage()
	allImg := map[int]map[int]image.Image{}

	count := 9
	for i := 1; i <= 4; i++ {
		if i == 4 {
			count = 8
		}
		for j := 1; j <= count; j++ {
			img := PickCardImage(newCard, j, i)
			imgs, ok := allImg[i]
			if !ok {
				imgs = map[int]image.Image{}
			}
			imgs[j] = img
			allImg[i] = imgs
		}
	}
	return allImg
}

// CreatePNG 產生PNG圖片
func CreatePNG(filename string, img image.Image) {
	dir := filepath.Dir(filename)
	_, err := os.Lstat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic(err)
		}
	}

	cf, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = png.Encode(cf, img)
	if err != nil {
		panic(err)
	}
	cf.Close()
}

// CreateGoFile 產生Go檔案
func CreateGoFile(filename, varname string, img image.Image) {
	dir := filepath.Dir(filename)
	_, err := os.Lstat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic(err)
		}
	}

	bb := bytes.NewBuffer(nil)
	err = png.Encode(bb, img)
	if err != nil {
		panic(err)
	}

	cf, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(cf, `package majo

var %s = []byte(%q)
`, varname, bb.Bytes())
	if err != nil {
		panic(err)
	}
	cf.Close()
}

// PickCardImage 挑取麻將卡片影像
func PickCardImage(card image.Image, ax, ay int) image.Image {
	ax--
	ay--
	rgba := image.NewRGBA(image.Rect(0, 0, dx, dy))
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			rx := ax*dx + x
			ry := ay*dy + y
			rgba.Set(x, y, card.At(rx, ry))
		}
	}

	return rgba
}

// CreateCardPNG 產生麻將圖片PNG
func CreateCardPNG() {
	CreatePNG("card_new.png", CardImage())
}

// CardImage 麻將圖檔
func CardImage() image.Image {
	img := RawCardImage()

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	newCard := image.NewRGBA(image.Rect(0, 0, width, height))
	cover := CoverImage()
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if j >= dy*3+8 && j < dy*4 && i >= dx*7 && i < dx*8-5 {
				newCard.Set(i, j, cover.At(i%dx, j%dy))
			} else {
				newCard.Set(i, j, img.At(i, j))
			}
		}
	}

	return newCard
}

// CreateRawCardImage 產生麻將原檔
func CreateRawCardImage() {
	img := RawCardImage()
	CreatePNG("card.png", img)
}

// RawCardImage 麻將原檔
func RawCardImage() image.Image {
	r := bytes.NewReader(cardJPG)
	img, err := jpeg.Decode(r)
	if err != nil {
		panic(err)
	}
	return img
}

// CreateCoverPNG 產生封面圖片PNG
func CreateCoverPNG() {
	CreatePNG("cover.png", CoverImage())
}

// CoverImage 封面圖片
func CoverImage() *image.RGBA {
	cover := image.NewRGBA(image.Rect(0, 0, dx, dy))
	colorGreen := color.RGBA{
		R: 89,
		G: 191,
		B: 115,
		A: 255,
	}
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			if y >= 8 && x < dx-5 {
				cover.Set(x, y, colorGreen)
			} else {
				cover.Set(x, y, color.White)
			}
		}
	}

	return cover
}
