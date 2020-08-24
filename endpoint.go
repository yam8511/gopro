package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var wg = sync.WaitGroup{}

// MessageEndpoint 處理訊息
func MessageEndpoint(bot *tgbotapi.BotAPI, uid int, msg *tgbotapi.Message) {
	if msg.IsCommand() {
		fmt.Printf("[#%d][%s]收到指令, 來自 %s\n指令: %s\n參數: %s\n附上: %s\n",
			msg.MessageID,
			msg.Time().Format("2006-01-02 15:04:05-0700"),
			msg.From.String(),
			msg.Command(),
			msg.CommandArguments(),
			msg.CommandWithAt(),
		)

		rply := NewReplyMessage(msg.Chat.ID, msg.MessageID, "")
		switch cmd := msg.Command(); cmd {
		case "help", "h":
			rply.Text = `可以輸入 /hi, /myid, 
			/keyboard, /closekeyboard, /inlineKeyboard`
		case "inlineKeyboard":
			rply.Text = "傳送inline鍵盤給你囉～"
			rply.ReplyMarkup = textKeyboard
		case "keyboard":
			rply.Text = "傳送數字鍵盤給你囉～"
			rply.ReplyMarkup = numericKeyboard
		case "closekeyboard":
			rply.Text = "關閉鍵盤"
			rply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "hi":
			rply.Text = "歡迎 telegram ~ " + msg.From.String()
		case "myid":
			rply.Text = fmt.Sprintf("你的ID是(%d)", msg.From.ID)
		default:
			rply.Text = fmt.Sprintf("未定義指令 [%s], 說明可以輸入 /help", cmd)
		}
		bot.Send(rply)
		return
	}

	fmt.Printf("[#%d][%s]收到新訊息, 來自 %s\n訊息內容: %s\n",
		msg.MessageID,
		msg.Time().Format("2006-01-02 15:04:05-0700"),
		msg.From.String(),
		msg.Text,
	)

	if msg.Sticker != nil {
		sticker := tgbotapi.NewStickerShare(msg.Chat.ID, msg.Sticker.FileID)
		bot.Send(sticker)
		return
	}

	rply := NewReplyMessage(msg.Chat.ID, msg.MessageID, "copy that")
	bot.Send(rply)
}

// CallbackEndpoint 反饋
func CallbackEndpoint(bot *tgbotapi.BotAPI, uid int, cb *tgbotapi.CallbackQuery) {
	fmt.Printf("[#%d][%s]收到新訊息, 來自 %s\n訊息內容: %s\n回傳資料: %s\nChat Instance: %s\nShort Name: %s\ncb ID: %s\ninline message id: %s\n",
		cb.Message.MessageID,
		cb.Message.Time().Format("2006-01-02 15:04:05"),
		cb.From.String(),
		cb.Message.Text,
		cb.Data,
		cb.ChatInstance,
		cb.GameShortName,
		cb.ID,
		cb.InlineMessageID,
	)

	if cb.Data == "more" {
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("新的回應🥶", "回傳Data按鈕"),
				tgbotapi.NewInlineKeyboardButtonSwitch("歡迎大家分享🥺", "回傳Switch按鈕"),
				tgbotapi.NewInlineKeyboardButtonURL("訂閱並且分享按讚🥳", "https://www.youtube.com/watch?v=FsFZOtDowF0&t=70s"),
			),
		)

		reply := tgbotapi.NewEditMessageText(
			cb.Message.Chat.ID,
			cb.Message.MessageID,
			"升級選項",
		)
		reply.ReplyMarkup = &markup
		bot.Send(reply)
	} else {
		ca := tgbotapi.NewCallback(cb.ID, "跳出CB")
		res, err := bot.AnswerCallbackQuery(ca)
		if err != nil || !res.Ok {
			ca = tgbotapi.NewCallbackWithAlert(cb.ID, "CB錯誤: "+err.Error())
			bot.AnswerCallbackQuery(ca)
		}
	}
}

// EditMessageEndpoint 處理訊息
func EditMessageEndpoint(bot *tgbotapi.BotAPI, uid int, msg *tgbotapi.Message) {
	fmt.Printf("[#%d][%s]有新的訊息, 來自 %s\n訊息內容: %s\n",
		msg.MessageID,
		msg.Time().Format("2006-01-02 15:04:05-0700"),
		msg.From.String(),
		msg.Text,
	)
	rply := NewReplyMessage(msg.Chat.ID, msg.MessageID, "catcha that")
	bot.Send(rply)
}

// UpdateMaster 更新入口
func UpdateMaster(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil {
		go func(update tgbotapi.Update) {
			wg.Add(1)
			defer wg.Done()
			defer func(uid int) {
				if r := recover(); r != nil {
					log.Printf("[#%d]發生Panic! 原因[%v]\n", uid, r)
				}
			}(update.UpdateID)
			MessageEndpoint(bot, update.UpdateID, update.Message)
		}(update)
	}

	if update.CallbackQuery != nil {
		go func(update tgbotapi.Update) {
			wg.Add(1)
			defer wg.Done()
			defer func(uid int) {
				if r := recover(); r != nil {
					log.Printf("[#%d]發生Panic! 原因[%v]\n", uid, r)
					debug.PrintStack()
				}
			}(update.UpdateID)
			CallbackEndpoint(bot, update.UpdateID, update.CallbackQuery)
		}(update)
	}

	if update.EditedMessage != nil {
		go func(update tgbotapi.Update) {
			wg.Add(1)
			defer wg.Done()
			defer func(uid int) {
				if r := recover(); r != nil {
					log.Printf("[#%d]發生Panic! 原因[%v]\n", uid, r)
				}
			}(update.UpdateID)
			EditMessageEndpoint(bot, update.UpdateID, update.EditedMessage)
		}(update)
	}
}

// WaitShutDown 等待程結束
func WaitShutDown(updateChan tgbotapi.UpdatesChannel) {
	done := make(chan byte)
	go func() {
		updateChan.Clear()
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("執行緒全部完成，結束程序")
	case <-time.After(time.Second * 30):
		log.Println("超過30秒，強制結束程序")
	}
}
