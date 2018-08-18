package main

import (
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

	b.Handle(tb.OnText, func(m *tb.Message) {
		b.Send(m.Sender, m.Text+" , 接收訊息 By "+m.Sender.Username)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		log.Printf("檔案資訊\n%+#v\n", m.Photo)
		log.Printf("檔案資訊\n%+#v\n", m.Photo.MediaFile())
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

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello~ "+m.Sender.Username)
	})

	b.Handle("/cmd", func(m *tb.Message) {
		b.Send(m.Sender, "你輸入了 "+m.Payload)
	})

	// This button will be displayed in user's
	// reply keyboard.
	replyBtn := tb.ReplyButton{Text: "🌕 Button #1"}
	replyBtn2 := tb.ReplyButton{Text: "🌕 Button #2"}
	replyBtn3 := tb.ReplyButton{Text: "🌕 Button #Contact", Contact: true}
	replyBtn4 := tb.ReplyButton{Text: "🌕 Button #Location", Location: true}
	replyBtn5 := tb.ReplyButton{Text: "🌕 Button #5", Action: func(cb *tb.Callback) {
		log.Printf("BTN5 ---> %+v\n", cb)
	}}
	replyKeys := [][]tb.ReplyButton{
		[]tb.ReplyButton{replyBtn, replyBtn2},
		[]tb.ReplyButton{replyBtn3, replyBtn4, replyBtn5},
		// ...
	}

	replyCB := func(m *tb.Message) {
		// on reply button pressed
		log.Println("reply key press -->", m.Text)
	}

	replyLocation := func(m *tb.Message) {
		// on reply button pressed
		// log.Println("reply key press -->", m.Text)
		// log.Printf("Location: %+v\n\n, Contact : %+v\n\n", m.Location, m.Contact)
		log.Printf("Location: %+v\n\n", m)
	}

	replyContact := func(m *tb.Message) {
		// on reply button pressed
		// log.Println("reply key press -->", m.Text)
		// log.Printf("Location: %+v\n\n, Contact : %+v\n\n", m.Location, m.Contact)
		log.Printf("Contact: %+v\n\n", m)
	}
	b.Handle(&replyBtn, replyCB)
	b.Handle(&replyBtn2, replyCB)
	b.Handle(&replyBtn3, replyContact)
	b.Handle(&replyBtn4, replyLocation)
	b.Handle(&replyBtn5, replyCB)

	// And this one — just under the message itself.
	// Pressing it will cause the client to send
	// the bot a callback.
	//
	// Make sure Unique stays unique as it has to be
	// for callback routing to work.
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

		log.Println(replyKeys)

		log.Println("回傳鍵盤")
		b.Send(m.Sender, "Hello!", &tb.ReplyMarkup{
			ReplyKeyboard:  replyKeys,
			InlineKeyboard: inlineKeys,
		})
	})

	b.Handle("/cancel", func(m *tb.Message) {
		b.Send(m.Sender, "Reset", &tb.ReplyMarkup{
			ReplyKeyboardRemove: true,
			// InlineKeyboard: inlineKeys,
		})
	})

	b.Start()
}
