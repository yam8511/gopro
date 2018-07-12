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

	// Demo 2
	c, d := demo2()
	fmt.Println(c, d)

	x, y := demo()
	fmt.Println(x, y)

	/**
	 * 時間用法
	 */

	// 現在時間
	// os.Setenv("TZ", "Asia/Taipei")
	now := time.Now()
	fmt.Println(now)

	// 時間轉時戳
	ts := now.Unix()
	fmt.Println(ts)
	// 時戳轉時間
	backts := time.Unix(ts, 0)
	fmt.Println(backts)
	// php Y-m-d H:i:s
	// 時間格式化
	fmt.Println("現在時間：", now.Format("2006-01-02 15:00:00"))

	// 字串轉時戳
	t := "2018-02-19 12:30:13+0800"
	tt, err := time.Parse("2006-01-02 15:04:05-0700", t)
	fmt.Println(tt.Format(time.RFC3339), err)

	// 睡覺
	// time.Sleep(time.Second * 3)
	// // 奈秒 time.Nanosecond
	// // 微秒 time.Microsecond
	// // 毫秒 time.Millisecond
	// // 秒   time.Second
	// // 分   time.Minute
	// // 時   time.Hour

	// 經過時間
	excursionn := time.Since(now)
	fmt.Println(excursionn)

	// 等待多少時間
	ch := time.After(time.Second * 30)
	apiCh := time.After(time.Nanosecond)
	fmt.Println("Start Sleep")
	select {
	case <-ch:
		fmt.Println("Timeout")
	case t := <-apiCh:
		fmt.Println(t.Sub(now))
		fmt.Println("Done")
	}
	fmt.Println("Done2")

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
	fmt.Println(strings.Contains("abcd", "bd"))
	// 是否包含任何字元
	fmt.Println(strings.ContainsAny("abc", "bd"))
	// 轉換
	fmt.Println(strings.Replace("abbscd", "b", "z", 3))

	/**
	 * 字串轉型處理
	 */
	// 字串轉數字
	numStr := "12"
	num, err := strconv.Atoi(numStr)
	fmt.Println(num, err)
	num64, err := strconv.ParseInt(numStr, 0, 32)
	fmt.Println(int32(num64), err)

	// 字串轉布林
	boolStr := "true"
	boolean, err := strconv.ParseBool(boolStr)
	fmt.Println(boolean, err)

	/**
	 * 併發、讀寫衝突、讀寫鎖
	 */
	// mx := new(sync.RWMutex)
	m := map[int]int{}
	for i := 0; i < 100; i++ {
		m[i] = i
	}
}

// 只宣告型態當回傳
func demo() (string, bool) {
	return "demo", false
}

// 以變數當回傳
func demo2() (x string, y bool) {
	x = "OK"
	y = true
	return "NO", false
}
