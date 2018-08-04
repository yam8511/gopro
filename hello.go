package main

import "fmt"

type demo struct {
}

func (demo) GetMe() {
}

type demo2 struct {
}

func (*demo2) GetMe() {
}

func main() {
	// 1. 定義規則
	// 2. 實現規則
	// 3. 執行結果

	// 可為任意值
	var data interface {
		GetMe()
	}
	fmt.Println("data --->", data)

	// 有實現 GetMe 的 func
	newDemo := demo{}
	data = newDemo
	fmt.Println("data --->", data)

	// 會噴錯，實現 GetMe 的 func 不能有 pointer
	newDemo2 := demo2{}
	newDemo2.GetMe()
	data = newDemo2

	// 會噴錯，int 並沒有實現 GetMe
	data = 1
	fmt.Println("data --->", data)

	// 會噴錯，string 並沒有實現 GetMe
	data = "OK"
	fmt.Println("data --->", data)
}
