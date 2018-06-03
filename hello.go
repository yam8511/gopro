package main

import (
	"fmt"
)

type demo struct {
	num int
}

func main() {
	s := []demo{demo{num: 1}, demo{num: 2}, demo{num: 3}, demo{num: 4}, demo{num: 5}}
	m := map[string]demo{
		"a": demo{num: 1},
		"b": demo{num: 2},
		"c": demo{num: 3},
	}
	fmt.Printf("slice -> %p\n", &s)
	fmt.Printf("map -> %p\n", &m)
	for i, d := range s {
		fmt.Printf("before, s[%d] -> %p\n", i, &s[i])
		fmt.Printf("before, %d -> %p\n", i, &d)
		d.num *= 2
		s[i] = d
		fmt.Printf("after, s[%d] -> %p\n", i, &s[i])
		fmt.Printf("after, %d -> %p\n", i, &d)
	}
	for i, d := range m {
		fmt.Printf("before, m[%s] -> %p\n", i, &d)
		d.num *= 3
		m[i] = d
		fmt.Printf("after, m[%s] -> %p\n", i, &d)
	}
	fmt.Printf("slice -> %p\n", &s)
	fmt.Printf("map -> %p\n", &m)
	for i, d := range s {
		fmt.Println("s -> ", i, d.num)
	}
	for i, d := range m {
		fmt.Println("m -> ", i, d.num)
	}
	fmt.Println("Hello World")
}
