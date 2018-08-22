package main

import (
	"fmt"
	"net"
)

func main() {
	checkIPInfo("192.168.1.2")
	checkIPInfo("FE80::0202:B3FF:FE1E:8329")
	checkIPInfo("192.168.i.2")
	checkIPInfo("FE80:0202:B3FF:FE1E:8329")
}

func checkIPInfo(address string) {
	// 解析IP
	ip := net.ParseIP(address)
	// 檢查IP是否有效
	if ip.To16() == nil {
		fmt.Println(address, "is not a valid IP address")
	} else {
		fmt.Println(address, "to IPv4  --->", ip.To4())
		fmt.Println(address, "to IPv16 --->", ip.To16())
	}

	/**
	 * ============
	 *     心得
	 * ============
	 * net.IP 有 To4() 與 To16()
	 * To4() 可以取得 IPv4 的格式
	 * To16() 可以取得 IPv6 的格式
	 *
	 * 然而使用 To16()，而非 To4()去判斷有沒有效
	 * 是因為 IPv4 可用 To16() 顯示，但 IPv6 無法用 To4()顯示
	 * 所以選擇 To16() 可涵蓋所有IP
	 *
	 */
}
