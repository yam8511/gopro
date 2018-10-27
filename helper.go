package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Welcome 介紹文字
func Welcome(nickname, username, person string) (welcome string) {
	welcome = fmt.Sprintf("您好~ 我是 %s (%s), 是一位 [%s]",
		nickname,
		username,
		person,
	)
	return
}

// GracefulDown 優雅結束程式
func GracefulDown() (sig chan os.Signal) {
	sig = make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)
	return
}
