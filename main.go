package main

import (
	"fmt"
	"strconv"
	"syscall/js"
)

func getElementByID(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

func add(i []js.Value) {
	value1 := getElementByID(i[0].String()).Get("value").String()
	value2 := getElementByID(i[1].String()).Get("value").String()

	int1, _ := strconv.Atoi(value1)
	int2, _ := strconv.Atoi(value2)

	// sum := i[0].Int() + i[1].Int()
	// js.Global().Set("output", js.ValueOf(sum))
	// println(js.ValueOf(sum).String())
	output := getElementByID(i[2].String())
	if output != js.Null() {
		output.Set("innerText", int1+int2)
	}
	js.Global().Set("output", int1+int2)
	println(int1 + int2)
}

func substract(i []js.Value) {
	int1 := i[0].Int()
	int2 := i[1].Int()
	i[2].Invoke(int1 - int2)
	js.Global().Set("output", int1-int2)
	println(int1 - int2)
}

// 方法中心
func goCall(args []js.Value) {
	if len(args) == 0 {
		fmt.Println("缺少參數")
		return
	}
	if args[0].Type() != js.TypeString {
		fmt.Println("方法名稱錯誤")
		return
	}
	method := args[0].String()
	fmt.Println("OK", method)
	if len(args) > 1 {
	}
	args = args[1:]
	fmt.Println("參數", args, len(args))
	cb, ok := methods[method]
	if !ok {
		fmt.Println("方法名稱錯誤")
		return
	}

	cb(args)
}

var methods = map[string]func([]js.Value){
	"Add":       add,
	"Substract": substract,
}

// 註冊給JS頁面可使用的方法名稱
func registerCallbacks() {
	for method, cb := range methods {
		js.Global().Set(method, js.NewCallback(cb))
	}
	js.Global().Set("GoCall", js.NewCallback(goCall))
}

func main() {
	registerCallbacks()
	println("Go WebAssembly Initialized.")
	c := make(chan int)
	<-c
}
