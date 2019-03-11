package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("請輸入時戳, 例如: ts 1552036736 1552037432")
		return
	}

	for i := range os.Args {
		if i == 0 {
			continue
		}
		tss := os.Args[i]
		ts, err := strconv.ParseInt(tss, 10, 64)
		if err != nil {
			fmt.Printf("%d) %s : NaN\n", i, tss)
			continue
		}

		t := time.Unix(ts, 0)
		fmt.Printf("%d) %s : %s\n", i, tss, t.Format(time.RFC3339))
	}
}
