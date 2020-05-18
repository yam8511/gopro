package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Up     = "w"
	Left   = "a"
	Right  = "d"
	Down   = "s"
	Adjust = "r"
	Finish = "c"
)

// 翻譯
func Trans(arrow string) string {
	switch arrow {
	case Up:
		return "前進"
	case Left:
		return "左轉"
	case Right:
		return "右轉"
	case Down:
		return "後退"
	default:
		return ""
	}
}

var trans = map[string]string{
	Up:    "前進",
	Left:  "左轉",
	Right: "右轉",
	Down:  "後退",
}

var ForwardMeterUnit, BackwardMeterUnit, RightAngleUnit, LeftAngleUnit time.Duration

var mapUnit = map[string]*time.Duration{
	Up:    &ForwardMeterUnit,
	Down:  &BackwardMeterUnit,
	Left:  &LeftAngleUnit,
	Right: &RightAngleUnit,
}

func main() {

	var moves []string
	var c string

	// 先校準單位
ADJUST:
	// 	adjust := func(action string, start, end func(), unit *time.Duration) {
	// 		log.Printf("倒數3秒後，將%s車子進行單位校準，移動一個單位後請按下Enter\n", action)
	// 		for i := 3; i > 0; i-- {
	// 			log.Println(i, ".. ")
	// 			time.Sleep(time.Second)
	// 		}

	// 		fmt.Println("請輸入任意鍵...")
	// 		now := time.Now()
	// 		start()
	// 		fmt.Scanln(&c)
	// 		*unit = time.Since(now)
	// 		end()
	// 	}

	// 	adjust("前進", func() {}, func() {}, &ForwardMeterUnit)
	// 	log.Println("前進距離平均一單是 ", ForwardMeterUnit)

	// 	adjust("右旋", func() {}, func() {}, &RightAngleUnit)
	// 	log.Println("右旋角度平均一單是 ", RightAngleUnit)

	// 	adjust("後退", func() {}, func() {}, &BackwardMeterUnit)
	// 	log.Println("後退距離平均一單是 ", BackwardMeterUnit)

	// 	adjust("左旋", func() {}, func() {}, &LeftAngleUnit)
	// 	log.Println("左旋角度平均一單是 ", LeftAngleUnit)

	for { // 非同步規劃路線
		log.Println("請輸入移動路徑(WASD)，完成請輸入C，重新校准請輸入R")
		moves = []string{}
	PLAN:
		for {
			_, err := fmt.Fscanf(os.Stdin, "%s", &c)
			if err != nil {
				continue
			}

			c = strings.ToLower(c)
			if c == Adjust {
				goto ADJUST
			}

			ss := strings.Split(c, "")
			for _, move := range ss {
				if move == Finish {
					break PLAN
				}

				if move == Up || move == Left || move == Right || move == Down {
					moves = append(moves, move)
				}
			}
		}

		run(moves)
	}
}

func run(moves []string) {
	if len(moves) == 0 {
		return
	}

	var forwardMoves, backMoves []string
	forwardMoves = append(forwardMoves, moves...)

	log.Println("路徑規劃 -> ", forwardMoves)
	for i := 0; i < len(forwardMoves); i++ {
		move := forwardMoves[i]
		fmt.Print(Trans(move) + " -> ")
	}
	log.Println()

	var firstMove = true
	for i := len(moves) - 1; i >= 0; i-- {
		move := moves[i]

		if firstMove {
			backMoves = append(backMoves, Right, Right)
		}

		switch move {
		case Up:
			backMoves = append(backMoves, Right, Right, Up)
		case Down:
			backMoves = append(backMoves, Up)
		case Left:
			backMoves = append(backMoves, Right)
		case Right:
			backMoves = append(backMoves, Left)
		}
	}

	log.Println("準備回去 -> ", backMoves)
	for i := 0; i < len(backMoves); i++ {
		move := backMoves[i]
		fmt.Print(Trans(move) + " -> ")
	}
	log.Println()
}

func GetUnit(action string) time.Duration {
	switch action {
	case Up:
		return ForwardMeterUnit
	case Down:
		return BackwardMeterUnit
	case Left:
		return LeftAngleUnit
	case Right:
		return RightAngleUnit
	default:
		return 0
	}
}
