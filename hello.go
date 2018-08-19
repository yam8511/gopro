package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	// [參考網址](https://github.com/tucnak/telebot)

	settings := tb.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tb.NewBot(settings)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/myid", func(m *tb.Message) {
		b.Send(m.Sender, "Your ID is "+fmt.Sprint(m.Sender.ID))
	})

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello~ "+m.Sender.Username)
	})

	b.Handle("/cmd", func(m *tb.Message) {
		b.Send(m.Sender, "You enter "+m.Payload)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		b.Send(m.Sender, m.Text+" , 接收訊息 By "+m.Sender.Username)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		log.Printf("檔案資訊\n%+#v\n", m.Photo)
		log.Printf("雲端？%+#v\n, 本地？%+#v", m.Photo.InCloud(), m.Photo.OnDisk())
		log.Printf("文字%s\n", m.Caption)
		// b.Send(m.Sender, m.Photo.FilePath+" , 接收檔案 By "+m.Sender.Username)
		b.Send(m.Sender, m.Photo)
	})

	b.Handle(tb.OnDocument, func(m *tb.Message) {
		log.Printf("檔案資訊\n%+#v\n", m.Document)
		log.Printf("雲端？%+#v\n, 本地？%+#v", m.Document.InCloud(), m.Document.OnDisk())
		log.Printf("文字%s\n", m.Caption)
		// b.Send(m.Sender, m.Photo.FilePath+" , 接收檔案 By "+m.Sender.Username)
		b.Send(m.Sender, m.Document)
	})

	// 以下設定額外的按鈕鍵盤
	replyBtn := tb.ReplyButton{Text: "🌕 Button #1"}
	replyBtn2 := tb.ReplyButton{Text: "🌕 Button #2"}
	replyBtn3 := tb.ReplyButton{Text: "🌕 Button #Contact", Contact: true}
	replyBtn4 := tb.ReplyButton{Text: "🌕 Button #Location", Location: true}
	replyBtn5 := tb.ReplyButton{Text: "🌕 Button #5"}
	replyKeys := [][]tb.ReplyButton{
		[]tb.ReplyButton{replyBtn, replyBtn2},
		[]tb.ReplyButton{replyBtn3, replyBtn4, replyBtn5},
		// ...
	}

	// 設定按鈕相對應的動作
	replyCB := func(m *tb.Message) {
		// on reply button pressed
		log.Println("reply key press -->", m.Text)
	}

	b.Handle(&replyBtn, replyCB)
	b.Handle(&replyBtn2, replyCB)
	b.Handle(&replyBtn3, replyCB)
	b.Handle(&replyBtn4, replyCB)
	b.Handle(&replyBtn5, replyCB)

	// 設定文字底下顯示的按鈕
	inlineBtn := tb.InlineButton{
		Unique:      "sad_moon",
		Text:        "🌚 Button #2",
		URL:         "https://www.youtube.com/watch?v=7iDPl1oLRNc",
		Data:        "zzz",
		InlineQuery: "what_query???",
	}
	inlineBtn2 := tb.InlineButton{
		Unique:      "happy_moon",
		Text:        "🌚 Button #3",
		Data:        "ccc",
		InlineQuery: "THIS_query???",
	}
	inlineKeys := [][]tb.InlineButton{
		[]tb.InlineButton{inlineBtn, inlineBtn2},
		// ...
	}

	inlineCB := func(c *tb.Callback) {
		// on inline button pressed (callback!)
		log.Println("inline key press -->", c.Data, c.ID)
		log.Printf("Message ---> %+v\n", c.Message)

		// always respond!
		b.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "COP", ShowAlert: true, URL: "https://github.com/"})
	}
	b.Handle(&inlineBtn, inlineCB)
	b.Handle(&inlineBtn2, inlineCB)

	// Command: /start <PAYLOAD>
	b.Handle("/start", func(m *tb.Message) {
		log.Println("私人???", m.Private())
		if !m.Private() {
			return
		}
		log.Println("回傳鍵盤")
		b.Send(m.Sender, "Hello!", &tb.ReplyMarkup{
			ReplyKeyboard:  replyKeys,
			InlineKeyboard: inlineKeys,
		})
	})

	b.Handle("/cancel", func(m *tb.Message) {
		b.Send(m.Sender, "Reset", &tb.ReplyMarkup{
			ReplyKeyboardRemove: true,
		})
	})

	b.Start()
}
