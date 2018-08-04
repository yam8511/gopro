package main

import (
	"log"
	"time"
)

func main() {
	// 一堆礦石
	theMine := []string{"rock", "ore", "ore", "rock", "rock"}
	// 建立兩個通道，一個是傳送找到的礦石，一個是傳送挖到的礦石
	foundOreChannel := make(chan string, 5)
	minedOreChannel := make(chan string, 5)
	// minedOreChannel := make(chan string)

	// 找礦石
	// foundOre := finder(theMine)
	// fmt.Println("找到了礦物 --->", foundOre)
	go finder(theMine, foundOreChannel)

	// 挖礦石
	// minedOre := miner(foundOre)
	// fmt.Println("挖好的礦物 --->", minedOre)
	go miner(foundOreChannel, minedOreChannel)

	// 提煉黃金
	// gold := smelter(minedOre)
	// fmt.Println("煉金好的黃金 --->", gold)
	smelter(minedOreChannel)
}

// 找礦者
func finder(mines []string, foundOreChannel chan string) {
	i := 0
	for _, mine := range mines {
		if mine == "ore" {
			time.Sleep(time.Second)
			log.Println("找到一個礦物 --->", i+1)
			// oreMines = append(oreMines, mine)
			foundOreChannel <- mine
			i++
		}
	}
	// 關閉通道
	log.Println("找礦者任務結束，關閉通道")
	close(foundOreChannel)
	return
}

// 挖礦者
func miner(foundOreChannel, minedOreChannel chan string) {
	i := 0
	for {
		mine := <-foundOreChannel
		if mine != "ore" {
			close(minedOreChannel)
			break
		}
		log.Println("開始挖一個礦物 --->", i+1)
		time.Sleep(time.Second)
		// oreMines = append(oreMines, mine)
		minedOreChannel <- mine
		log.Println("挖好了一個礦物 --->", i+1)
		i++
	}
	// return
}

// 煉金者
func smelter(minedOreChannel chan string) {
	i := 0
	// for i, mine := range mines {
	for {
		mine := <-minedOreChannel
		if mine != "ore" {
			return
		}
		log.Println("開始提煉一個黃金 --->", i+1)
		time.Sleep(time.Millisecond)
		// gold = append(gold, "黃金"+mine)
		log.Println("提煉好了一個黃金 --->", i+1)
	}
	// return
}
