package main

import "log"

func main() {
	a := []int{5, 4, 0, 3, 1, 6, 2}
	log.Println(a)
	tmp := map[int]int{}
	ans := 0

	index := 0
LOOP:
	for {
		if index >= len(a) { // 如果大於了，則直接紀錄
			ans++                 // 如果存在了，+1
			for k, n := range a { // 找過要找新的index
				_, ok := tmp[n]
				if !ok {
					index = k
					break
				} else if k == len(a)-1 { // 已經結束了
					break LOOP
				}
			}
			continue
		}

		v := a[index]

		_, ok := tmp[v]
		if ok {
			ans++                 // 如果存在了，+1
			for k, n := range a { // 找過要找新的index
				_, ok = tmp[n]
				if !ok {
					index = k
					break
				} else if k == len(a)-1 { // 已經結束了
					break LOOP
				}
			}
		} else {
			tmp[v] = index
			index = v
		}
	}

	log.Println(ans)
}
