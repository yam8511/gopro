package main

import (
	"fmt"
)

func main() {
	// 字串
	str := "ABC"
	// 字串轉[]Byte
	strToB := []byte(str)
	// []Byte轉字串
	bToStr := string(strToB)
	// 取字串中的一個字元
	oneOfStr := bToStr[0]
	// 硬出來
	fmt.Println(str, oneOfStr, string(oneOfStr))

	// 數字
	num := 123
	// 轉64位元的數字
	num64 := int64(num)
	fmt.Printf("num 的型態 %T, 數值 %d or %v, 位址 %p\n", num, num, num, &num)
	fmt.Printf("num64 的型態 %T, 數值 %d or %v, 位址 %p\n", num64, num64, num64, &num64)

	// 數字轉字串
	numToStr := fmt.Sprint(num64)
	fmt.Printf("numToStr 的型態 %T, 數值 %s or %v, 位址 %p\n", numToStr, numToStr, numToStr, &numToStr)

	// 小數點
	floatNum := 12.35
	fmt.Printf("floatNum 的型態 %T, 數值 %.1f or %v, 位址 %p\n", floatNum, floatNum, floatNum, &floatNum)

	/**
	 * fmt.Printf 的 f 就是 format
	 * %d, %2d, %3d, ... 整數
	 * %f, %.2f, %.3f, ... 浮點數
	 * %s 字串
	 * %T 型態
	 * %p 記憶體位址
	 * %v 不管型態，印出值
	 */

	// 數字陣列 (宣告變數時，即固定長度)
	numArray := [5]int{}
	fmt.Println("array", numArray)
	// 數字切片 (宣告變數時，無固定長度)
	numSlice := []int{}
	fmt.Println("slice", numSlice)
	// 數字切片 (宣告變數時，無固定長度，但一開始就會有5個元素)
	makeNumSlice := make([]int, 5)
	fmt.Println("make Slice", makeNumSlice)

	// Map
	aMap := map[string]int{
		"b": 2,
		"a": 1,
		"c": 3,
	}
	fmt.Println("map", aMap)
	bMap := make(map[string]int)
	fmt.Println("make map", bMap)

	// 關於迴圈，只有 for
	// 無窮迴圈
	for {
		fmt.Println("這是無窮迴圈")
		break
	}
	// 一般的迴圈 i++
	fmt.Println("===== 這是一般迴圈 ====")
	for i := 0; i < 10; i++ {
		fmt.Print(i)
	}
	fmt.Println()

	fmt.Println("===== For Map : Key -> Value ====")
	// php 的 foreach
	for key, value := range aMap {
		fmt.Printf("key %s, value %d\n", key, value)
	}

	// 只需要 key
	fmt.Println("===== For Map : Key ====")
	for key := range aMap {
		fmt.Printf("key %s, value %d\n", key, aMap[key])
	}

	// 只需要 value
	fmt.Println("===== For Map : Value ====")
	for _, value := range aMap {
		fmt.Printf("only value %d\n", value)
	}

	fmt.Println("===== For Array : Index -> Value ====")
	// php 的 foreach
	for index, value := range numSlice {
		fmt.Printf("index %d, value %d\n", index, value)
	}

	// 只需要 idnex
	fmt.Println("===== For Array : Index ====")
	for index := range numSlice {
		fmt.Printf("index %d, value %d\n", index, numSlice[index])
	}

	// 只需要 value
	fmt.Println("===== For Array : Value ====")
	for _, value := range numSlice {
		fmt.Printf("only value %d\n", value)
	}

	// 陣列塞入
	numSlice = append(numSlice, 1, 2, 3, 4)
	fmt.Println("===== For Array : Append ====")
	for _, value := range numSlice {
		fmt.Printf("only value %d\n", value)
	}

	// 先說說 ...
	takeout(1, 2, 3, 4)
	takeout(numSlice...)

	// 取陣列的特定範圍
	// [x:y] x 開啟，但不包含y
	// fmt.Println(numSlice[1:3])

	// 陣列取出第i個元素
	numSlice = []int{1, 2, 3, 4, 5, 6}
	i := 2
	numSlice = append(numSlice[:i], numSlice[i+1:]...)
	fmt.Println("===== For Array : Delete ====")
	for _, value := range numSlice {
		fmt.Printf("only value %d\n", value)
	}

	// php 的 isset
	aMap = map[string]int{
		"b": 2,
		"a": 1,
		"c": 3,
	}
	fmt.Println("===== For Map : isset ====")
	fmt.Println(aMap["a"])
	a1, ok := aMap["a"]
	fmt.Printf("a 值是 %d, 存在: %v\n", a1, ok)
	z1, ok := aMap["z"]
	fmt.Printf("z 值是 %d, 存在: %v\n", z1, ok)

	// 刪除 map 元素
	delete(aMap, "a")
	fmt.Println("===== For Map : Delete ====")
	for _, value := range aMap {
		fmt.Printf("only value %d\n", value)
	}

	// 提醒：switch 內建隱藏 break
	// 一般的 switch
	fmt.Println("===== For Switch ====")
	逼逼 := 1
	switch 逼逼 {
	case 0:
		fmt.Println("b is 0")
	case 1:
		fmt.Println("b is 1")
	default:
		fmt.Println("b is b")
	}

	// 特別的 switch (然後只會跑一個結果)
	switch {
	case 逼逼 == 0:
		fmt.Println("逼逼 is 0")
	case 逼逼 > 0:
		fmt.Println("逼逼 is > 0")
	case 逼逼 == 1:
		fmt.Println("逼逼 is 1")
	case 逼逼 < 3:
		fmt.Println("逼逼 is < 3")
	default:
		fmt.Println("逼逼 is b")
	}

	// 特別的 switch (讓他會跑多個結果)
	switch {
	case 逼逼 == 0:
		fmt.Println("逼逼 is 0")
		fallthrough
	case 逼逼 > 0:
		fmt.Println("逼逼 is > 0")
		fallthrough
	case 逼逼 == 4:
		fmt.Println("逼逼 is 1")
		fallthrough
	case 逼逼 < 3:
		fmt.Println("逼逼 is < 3")
		fallthrough
	default:
		fmt.Println("逼逼 is b")
	}

	var (
		都可以  interface{}
		只收數字 int
		// 只收字串 string
	)
	都可以 = "bbb"
	// 都可以 = 123
	fmt.Printf("都可以的形態是 %T, 值是 %v\n", 都可以, 都可以)

	只收數字, ok = 都可以.(int)
	fmt.Printf("轉成功嗎？%v ... 只收字串的形態是 %T, 值是 %v\n", ok, 只收數字, 只收數字)

	obj := Demo{"ZZ", 12}
	SetDemoName(obj)
	fmt.Println(obj.Name)
	SetDemoNameFromMemory(&obj)
	fmt.Println(obj.Name)
}

// SetDemoName 改名稱
func SetDemoName(d Demo) {
	d.Name = "SSS"
}

// SetDemoNameFromMemory 從記憶體改名稱
func SetDemoNameFromMemory(d *Demo) {
	d.Name = "SSS"
}

// Demo 這是範例
type Demo struct {
	Name string
	Age  int
}

// AAA 這是番宣
func (d *Demo) AAA(args string) string {
	return "I'm " + d.Name + args
}

// takeout ssss
func takeout(a ...int) {
	fmt.Println(a)
}
