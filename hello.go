package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	// ip := "192.168.1.*"
	ip := "FE80::0202:B3FF:FE1E:*"
	ips := parseIP(ip)
	for _, ip := range ips {
		checkIPInfo(ip)
	}
	checkIPInfo("192.168.1.2")
	checkIPInfo("FE80::0202:B3FF:FE1E:8329")
	checkIPInfo("192.168.i.2")
	checkIPInfo("FE80:0202:B3FF:FE1E:8329")
}

func parseIP(ip string) (ips []string) {
	ss := strings.Split(ip, ".")
	isIPv4 := false
	// var ips []string
	var ipElem [][]string
	if len(ss) == 4 {
		isIPv4 = true
		ipElem = [][]string{
			[]string{},
		}

		for i, s := range ss {
			switch {
			case strings.Contains(s, "*"):
				if i != 3 {
					return
				}
				for c := 0; c < 256; c++ {
					if c < len(ipElem) {
						ipElem[c] = append(ipElem[c], fmt.Sprint(c))
					} else {
						tmp := []string{ipElem[0][0], ipElem[0][1], ipElem[0][2], fmt.Sprint(c)}
						ipElem = append(ipElem, tmp)
					}
				}
			case strings.Contains(s, "/"):
				if i != 3 {
					return
				}
				elem := strings.Split(s, "/")
				if len(elem) != 2 {
					return
				}
				startIndex, err := strconv.Atoi(elem[0])
				if err != nil {
					return
				}
				endIndex, err := strconv.Atoi(elem[1])
				if err != nil {
					return
				}
				for c := startIndex; c <= endIndex; c++ {
					if c == startIndex {
						ipElem[0] = append(ipElem[0], fmt.Sprint(c))
					} else {
						tmp := []string{ipElem[0][0], ipElem[0][1], ipElem[0][2], fmt.Sprint(c)}
						ipElem = append(ipElem, tmp)
					}
				}
			default:
				for k := range ipElem {
					ipElem[k] = append(ipElem[k], s)
					ipElem[k][i] = s
				}
			}
		}
	} else {
		ss := strings.Split(ip, ":")
		fmt.Println(ss)
		if len(ss) == 6 {
			ipElem = [][]string{
				[]string{},
			}

			for i, s := range ss {
				switch {
				case strings.Contains(s, "*"):
					if i != 5 {
						return
					}
					for c := 0x0000; c < 0xFFFF; c += 0x0001 {
						elem := fmt.Sprintf("%X", c)
						if c < len(ipElem) {
							ipElem[c] = append(ipElem[c], elem)
						} else {
							tmp := []string{
								ipElem[0][0],
								ipElem[0][1],
								ipElem[0][2],
								ipElem[0][3],
								ipElem[0][4],
								elem,
							}
							ipElem = append(ipElem, tmp)
						}
					}
				case strings.Contains(s, "/"):
					if i != 5 {
						return
					}
					elem := strings.Split(s, "/")
					if len(elem) != 2 {
						return
					}

					startIndex, err := hex.DecodeString(elem[0])
					if err != nil {
						return
					}
					endIndex, err := hex.DecodeString(elem[1])
					if err != nil {
						return
					}

					c := []byte(string(startIndex))
					target := ByteAdd(endIndex, 0x0001)
					for !ByteEqual(c, target) {
						elem := fmt.Sprintf("%X", c)
						if ByteEqual(c, startIndex) {
							ipElem[0] = append(ipElem[0], elem)
						} else {
							tmp := []string{
								ipElem[0][0],
								ipElem[0][1],
								ipElem[0][2],
								ipElem[0][3],
								ipElem[0][4],
								elem,
							}
							ipElem = append(ipElem, tmp)
						}
						fmt.Printf("%X\n", c)
						ByteAdd(c, 0x01)
					}

				default:
					for k := range ipElem {
						ipElem[k] = append(ipElem[k], s)
						ipElem[k][i] = s
					}
				}
			}
		}
	}

	sep := ":"
	if isIPv4 {
		sep = "."
	}
	for _, elem := range ipElem {
		ips = append(ips, strings.Join(elem, sep))
	}
	return
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

// ByteEqual 檢查兩個[]byte是否相等
func ByteEqual(a, b []byte) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// ByteAdd []Byte加數字
func ByteAdd(c []byte, add byte) []byte {
	if c == nil {
		c = []byte{0}
		return c
	}
	index := len(c) - 1
	c[index] += add
	if c[index] == 0 {
		if index-1 < 0 {
			c[index]++
			c = append(c, 0)
		} else {
			c[index-1]++
		}
	}
	return c
}
