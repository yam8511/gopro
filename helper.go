package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/singleflight"
)

var pwd = ""
var programMX = &sync.RWMutex{}
var program *exec.Cmd
var fly = &singleflight.Group{}

// GetLocalIPs 取本機IP
func GetLocalIPs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return nil, err
	}

	ips := []string{}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.To4().String())
		}
	}
	return ips, nil
}

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
