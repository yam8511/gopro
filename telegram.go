package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var textKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("🚄", "https://irs.thsrc.com.tw/IMINT/"),
		tgbotapi.NewInlineKeyboardButtonSwitch("🔑", "key"),
		tgbotapi.NewInlineKeyboardButtonData("🔒", "lock"),
	),
)
var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

// NewBot 建立機器人
func NewBot(token string) (bot *tgbotapi.BotAPI, err error) {
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return
	}
	bot.Debug = true
	return
}

// CreateUpdateChannel 建立訊息更新通道
func CreateUpdateChannel(bot *tgbotapi.BotAPI, timeout, offset, limit int) (updateChan tgbotapi.UpdatesChannel, err error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout
	u.Offset = offset
	u.Limit = limit
	updateChan, err = bot.GetUpdatesChan(u)
	return
}

// NewReplyMessage 建立回覆訊息
func NewReplyMessage(chat int64, to int, text string) (msg tgbotapi.MessageConfig) {
	msg = tgbotapi.NewMessage(chat, text)
	msg.ReplyToMessageID = to
	return
}
