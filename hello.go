package main

import (
	"fmt"
	"net/url"
)

func main() {
	/**
	 * ============================
	 *      URL Query 初探心得
	 * ============================
	 * url.Values 本身型態是 map[string][]string
	 * Encode：顯示網址上的Query
	 *
	 * Add：可讓Query接上多個參數
	 * ------- ex. -------
	 * query 目前為 ""
	 * Add("a", "b") ---> "a=b"
	 * Add("a", "c") ---> "a=b&a=c"
	 *
	 * Set：直接重新設定指定欄位的參數值
	 * ------- ex. -------
	 * query 目前為 "a=b&a=c&a=d"
	 * Set("a", "z") ---> "a=z"
	 *
	 * Get：取指定欄位的第一個參數值
	 * ------- ex. -------
	 * query 目前為 "a=b&a=c&a=d"
	 * Get("a") ---> "a=b"
	 *
	 * Del：刪除指定欄位
	 * ------- ex. -------
	 * query 目前為 "a=b&a=c&x=y"
	 * Del("a") ---> "x=y"
	 *
	 *
	 * url.ParseQuery：用來解析網址上的Query
	 * 可將Query字串，轉成url.Values物件
	 *
	 */
	v := make(url.Values)
	v.Add("name", "zuolar")
	v.Add("name", "maius")
	v.Add("golang", "docker")
	v.Add("golang", "k8s")
	fmt.Println("Query --->", v.Encode())
	fmt.Println("Query Name --->", v.Get("name"))
	fmt.Println("Query Go --->", v.Get("golang"))

	v.Set("name", "zuolar")
	fmt.Println("Query --->", v.Encode())

	v.Del("golang")
	fmt.Println("Query --->", v.Encode())

	d, err := url.ParseQuery("golang=docker&golang=k8s&name=zuolar&#zzz")
	if err != nil {
		fmt.Println("Parse Query Error --->", err)
		return
	}
	fmt.Println("Parse Query", d)
}
