package main

import (
	"errors"
	"fmt"
)

/**
 * ===============================
 * 本分支大概說明 Error 的定義
 * ===============================
 * error 其實是一個 interface
 * 只要有實現 Error() string 的func
 * 就可以稱作 error 型態
 * error的定義如下
 *
 * type error interface {
 *     Error() string
 * }
 */

// 這邊定義一個新的struct，叫做 isZuolarCode
type isZuolarCode struct {
	Code string
}

// 給 isZuolarCode 賦予一個 func
// 叫做 Error() string  ---> 有符合 error 的定義
func (co isZuolarCode) Error() string {
	return fmt.Sprintf("Code [%s] is zuolar's code", co.Code)
}

// 檢查是否為 isZuolarCode 的錯誤
func isZuolarError(err error) bool {
	_, ok := err.(isZuolarCode)
	return ok
}

// 檢查是否為數字 ---> 這邊要對照 isZuolarError 的func，因為作法是一樣的
func isNumber(what interface{}) bool {
	_, ok := what.(int)
	return ok
}

func main() {
	fmt.Println("Hello World")

	// 簡單說明 interface{} 的轉型功能
	var anyData interface{}
	// 先賦予數字
	anyData = 1
	// 但是可以用 .(型態) 進行型態轉型
	// 會回傳兩個值，第一個是轉型後的資料，第二個是用來判斷是否轉型成功
	stringVar, convertOK := anyData.(string) // 先轉成字串
	// 如果轉型成功，顯示出來
	if convertOK {
		fmt.Println("字串轉型成功 --->", stringVar)
	} else {
		fmt.Println("字串轉型失敗 --->", anyData)
	}

	// 如果轉型失敗，改轉成數字
	numberVar, convertOK := anyData.(int)
	// 如果轉型成功，顯示出來
	if convertOK {
		fmt.Println("數字轉型成功 --->", numberVar)
	} else {
		fmt.Println("數字轉型失敗 --->", anyData)
	}
	fmt.Println("==========================")

	// 呼叫func來檢查
	fmt.Printf("anyData的變數值是不是數字？ %v\n", isNumber(anyData))
	fmt.Println("==========================")

	// golang原生基本的error
	normalErr := errors.New("一般的Error")
	fmt.Printf("錯誤型態: %T,\n 錯誤值: %v,\n 是否為 Zuolar 的錯誤？ %v\n", normalErr, normalErr, isZuolarError(normalErr))

	fmt.Println("==========================")
	// 自定義的error
	zuolarErr := isZuolarCode{Code: "Hello World"}
	fmt.Printf("錯誤型態: %T,\n 錯誤值: %v,\n 是否為 Zuolar 的錯誤？ %v\n", zuolarErr, zuolarErr, isZuolarError(zuolarErr))
}
