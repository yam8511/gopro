package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	/**
	 * func 宣告方式
	 */

	// Demo 1
	a, b := demo()
	fmt.Println(a, b)

	// Demo 2
	c, d := demo2()
	fmt.Println(c, d)

	/**
	 * 時間用法
	 */

	// 現在時間
	now := time.Now()
	fmt.Println(now)

	// 時間轉時戳
	ts := now.Unix()
	fmt.Println(ts)

	// 字串轉時戳
	t := "2018-02-19 12:30:13"
	tt, err := time.Parse("2006-01-02 15:04:05", t)
	fmt.Println(tt, err)

	// 睡覺
	time.Sleep(time.Nanosecond)
	// 奈秒 time.Nanosecond
	// 微秒 time.Microsecond
	// 毫秒 time.Millisecond
	// 秒   time.Second
	// 分   time.Minute
	// 時   time.Hour

	// 經過時間
	excursionn := time.Since(now)
	fmt.Println(excursionn)

	// 等待多少時間
	<-time.After(time.Millisecond)
	fmt.Println("Done")

	// 時間相減
	fmt.Println(time.Now().Sub(now))

	/**
	 * 字串處理
	 */

	// 去除前後空白
	fmt.Println(strings.TrimSpace(" aa "))
	// 去頭
	fmt.Println(strings.TrimPrefix("#aa", "#a"))
	// 去尾
	fmt.Println(strings.TrimSuffix("aa#", "#"))
	// 是否包含
	fmt.Println(strings.Contains("abc", "bd"))
	// 是否包含任何字元
	fmt.Println(strings.ContainsAny("abc", "bd"))

	/**
	 * 字串轉型處理
	 */

	// 字串轉數字
	numStr := "12"
	num, err := strconv.Atoi(numStr)
	fmt.Println(num, err)

	// 字串轉布林
	boolStr := "true"
	boolean, err := strconv.ParseBool(boolStr)
	fmt.Println(boolean, err)
}

// 只宣告型態當回傳
func demo() (string, bool) {
	return "demo", false
}

// 以變數當回傳
func demo2() (x string, y bool) {
	x = "OK"
	y = true
	return
}
