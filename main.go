package main

import (
	"log"
	"os"
	"time"
)

func main() {
	bot, err := NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("創建API機器人失敗,", err.Error())
	}

	bot.RemoveWebhook()

	person := "NPC"
	if !bot.Self.IsBot {
		person = "人類"
	}
	welcome := Welcome(bot.Self.String(), bot.Self.UserName, person)
	log.Println(welcome)

	log.Println("建立Telegram通道中 ...")
	updateChan, err := CreateUpdateChannel(bot, 60, 0, 0)
	if err != nil {
		log.Fatal("建立Telegram通道失敗,", err.Error())
	}
	log.Println("Telegram通道清除舊訊息中 ...")
	time.Sleep(time.Millisecond * 500)
	updateChan.Clear()
	log.Println("建立Telegram通道成功, 開始等待訊息")

	for {
		select {
		case sig := <-GracefulDown():
			log.Printf("接收到系統訊號[%s], 等待結束...\n", sig.String())
			WaitShutDown(updateChan)
			return
		case update := <-updateChan:
			go UpdateMaster(bot, update)
		}
	}
}
