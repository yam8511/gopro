package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/smtp"
	"strings"
)

// 以下 variable 可參考 Gmail 的 smtp 設定說明
var (
	host     = "smtp.gmail.com:587"
	username = "example@gmail.com"
	password = "passsword"
)

func main() {
	// SMTP設定
	auth := smtp.PlainAuth(host, username, password, "smtp.gmail.com")
	to := []string{"example@gmail.com"}
	boundary := "hhhrrr"

	// Like Mail Header
	mailBody := &bytes.Buffer{}
	mailBody.WriteString("Subject:測試信件\r\n")
	mailBody.WriteString("From: " + username + "\r\n")
	mailBody.WriteString("To: " + strings.Join(to, ",") + "\r\n")
	mailBody.WriteString(`Content-Type: multipart/mixed; boundary="` + boundary + `"` + "\r\n")

	// Mail Body Writer
	w := multipart.NewWriter(mailBody)
	defer w.Close()
	w.SetBoundary(boundary)

	// Mail Body
	bw, err := w.CreatePart(map[string][]string{})
	if err != nil {
		log.Println("建立信件內文失敗,", err.Error())
		return
	}

	bw.Write([]byte("你好，這裡有一張圖片。"))
	bw.Write([]byte("請記得打開來看！\r\n"))
	bw.Write([]byte("謝謝你\r\n"))

	// 附件
	imageByte, err := ioutil.ReadFile("go2.png")
	if err != nil {
		log.Println("加上附件, 讀檔失敗,", err.Error())
		return
	}

	fw, err := w.CreatePart(map[string][]string{
		"Content-Type": []string{`image/png`},
		"Content-Disposition": []string{
			`attachment; filename="go2.png"`,
		},
		"Content-Transfer-Encoding": []string{"base64"},
	})
	if err != nil {
		log.Println("建立信件附件內文失敗,", err.Error())
		return
	}
	encoder := base64.NewEncoder(base64.StdEncoding, fw)
	encoder.Write(imageByte)

	err = smtp.SendMail(
		host,
		auth,
		username,
		to,
		mailBody.Bytes(),
	)

	if err != nil {
		log.Println("寄送失敗, ", err.Error())
	}

}
