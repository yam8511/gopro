package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Photo 照片，以下的「photo」是我自取的標籤名稱
type Photo struct {
	Name      string    `photo:"name" pic:"pic_name"` // 有給 photo & pic 的標籤
	Size      int       `photo:"size" pic:"pic_size"` // 有給 photo & pic 的標籤
	CreatedAt time.Time `photo:"created_at"`          // 只有給 photo 的標籤
}

func main() {
	// 建立一個物件
	handsome := Photo{}

	// 取這個物件的結構
	handsomeType := reflect.TypeOf(handsome)

	// 顯示資料屬性
	fmt.Printf("Photo 這個物件有幾個欄位？ %d 個\n", handsomeType.NumField())
	fieldName := []string{}
	tagPhotoName := []string{}
	tagPicName := []string{}
	for i := 0; i < handsomeType.NumField(); i++ {
		// 欄位資料
		fieldInfo := handsomeType.Field(i)
		fieldName = append(fieldName, fieldInfo.Name)

		// 標籤資料
		tagInfo := fieldInfo.Tag

		// 尋找是否有「photo」的標籤
		valueInTag, hasPhotoTag := tagInfo.Lookup("photo")
		// 如果有「photo」的標籤資料，存進去變數
		if hasPhotoTag {
			tagPhotoName = append(tagPhotoName, valueInTag)
		} else {
			tagPhotoName = append(tagPhotoName, "<沒資料>")
		}

		// 尋找是否有「pic」的標籤
		valueInTag, hasPicTag := tagInfo.Lookup("pic")
		// 如果有「pic」的標籤資料，存進去變數
		if hasPicTag {
			tagPicName = append(tagPicName, valueInTag)
		} else {
			tagPicName = append(tagPicName, "<沒資料>")
		}

	}
	fmt.Printf("欄位名稱分別是 %s \n", strings.Join(fieldName, ", "))
	fmt.Printf("標籤「photo」的資料分別是 %s \n", strings.Join(tagPhotoName, ", "))
	fmt.Printf("標籤「pic」的資料分別是 %s \n", strings.Join(tagPicName, ", "))
}
