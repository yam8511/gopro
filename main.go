package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("請輸入時戳或時間格式, 例如: ts 1552036736 1552037432 2019-04-10T12:34:56Z")
		return
	}

	for i := range os.Args {
		if i == 0 {
			continue
		}
		tss := os.Args[i]
		ts, err := strconv.ParseInt(tss, 10, 64)
		if err != nil {
			t, err := time.Parse(time.RFC3339, tss)
			if err != nil {
				fmt.Printf("%d) %s : NaN\n", i, tss)
				continue
			}
			fmt.Printf("%d) %s : %d\n", i, tss, t.Unix())
			continue
		}

		t := time.Unix(ts, 0)
		fmt.Printf("%d) %s : %s\n", i, tss, t.Format(time.RFC3339))
	}
}
