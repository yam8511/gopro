package main

import (
	"fmt"
	"log"
	"time"

	lua "github.com/yuin/gopher-lua"

	starlight "github.com/starlight-go/starlight"
	"go.starlark.net/starlark"
)

func main() {
	fmt.Println("Hello World")
	count := 1

	fmt.Println("========")
	now := time.Now()
	for i := 0; i < count; i++ {
		light()
	}
	fmt.Println(time.Since(now))
	fmt.Println("========")
	now = time.Now()
	for i := 0; i < count; i++ {
		googleStarlark()
	}
	fmt.Println(time.Since(now))
	fmt.Println("========")
	now = time.Now()
	for i := 0; i < count; i++ {
		dolua()
	}
	fmt.Println(time.Since(now))
}

func dolua() {
	L := lua.NewState()
	L.SetGlobal("x", lua.LNumber(10))
	defer L.Close()
	if err := L.DoFile("fib.lua"); err != nil {
		panic(err)
	}
}

func light() {

	// errors will tell you about syntax/runtime errors.
	_, err := starlight.Eval("fibonacci.star", map[string]interface{}{
		"y": 0,
		"x": 10,
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	// for k, v := range values {
	// 	log.Printf("%s -> %T -> %v\n", k, v, v)
	// }
}

func googleStarlark() {

	// Execute Starlark program in a file.
	thread := &starlark.Thread{Name: "my thread"}

	pre := starlark.StringDict(map[string]starlark.Value{
		"x": starlark.MakeInt(10),
		"y": starlark.MakeInt(0),
	})
	_, err := starlark.ExecFile(thread, "fibonacci.star", nil, pre)
	if err != nil {
		log.Fatal("Exec ", err)
	}

	// for k, v := range values {
	// 	log.Printf("%s -> %T -> %v\n", k, v, v)
	// 	l := v.(*starlark.List)
	// 	for i := 0; i < l.Len(); i++ {
	// 		log.Println(l.Index(i).String())
	// 	}
	// }
}
