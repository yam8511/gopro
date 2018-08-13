package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {

	// 一般字串相加
	begin := time.Now()
	ss := "Hi there "
	for i := 0; i < 100000; i++ {
		ss += "(" + fmt.Sprint(i) + "," + fmt.Sprint(i) + "),"
	}
	fmt.Println("一般字串相加 ---> ", time.Since(begin))

	// 使用 string.Builder 做字串相加
	begin = time.Now()
	s := new(strings.Builder)
	s.WriteString("Hi there ")
	for i := 0; i < 100000; i++ {
		s.WriteString("(" + fmt.Sprint(i) + "," + fmt.Sprint(i) + "),")
	}
	fmt.Println("strings.Builder ---> ", time.Since(begin))

	// 驗證字串是否相同
	same := s.String() == ss
	fmt.Println("兩者組成的字串是否相同？", same)
}
