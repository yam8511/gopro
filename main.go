package main

import (
	"strconv"
	"syscall/js"
)

func add(i []js.Value) {
	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()
	value2 := js.Global().Get("document").Call("getElementById", i[1].String()).Get("value").String()

	int1, _ := strconv.Atoi(value1)
	int2, _ := strconv.Atoi(value2)

	// sum := i[0].Int() + i[1].Int()
	// js.Global().Set("output", js.ValueOf(sum))
	// println(js.ValueOf(sum).String())
	js.Global().Get("document").Call("getElementById", i[2].String()).Set("value", int1+int2)
	js.Global().Set("output", int1+int2)
	println(int1 + int2)
}

func substract(i []js.Value) {
	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()
	value2 := js.Global().Get("document").Call("getElementById", i[1].String()).Get("value").String()

	int1, _ := strconv.Atoi(value1)
	int2, _ := strconv.Atoi(value2)

	js.Global().Get("document").Call("getElementById", i[2].String()).Set("value", int1-int2)
	js.Global().Set("output", int1-int2)
	println(int1 - int2)
}

func registerCallbacks() {
	js.Global().Set("Add", js.NewCallback(add))
	js.Global().Set("Substract", js.NewCallback(substract))
}

func main() {
	c := make(chan int)

	println("Go WebAssembly Initialized.")
	registerCallbacks()

	<-c
}
