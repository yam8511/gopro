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

	b.Start()
}
