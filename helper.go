package main

import (
	"log"
	"net"
)

func isLocalIP(ip string) bool {
	_, ok := localIPs[ip]
	return ok
}

// 本機可用IP
func localAvailableIPs() map[string]string {
	availableIPs := map[string]string{}

	cfgs, err := net.Interfaces()
	if err != nil {
		log.Fatal("取本機取網卡資訊失敗: " + err.Error())
		return availableIPs
	}

	for _, cfg := range cfgs {
		ips, err := cfg.Addrs()
		if err != nil {
			log.Fatal("取本機取網卡資訊失敗: " + err.Error())
			continue
		}

		for _, ip := range ips {
			iip, _, err := net.ParseCIDR(ip.String())
			if err != nil {
				log.Fatal("取本機取網卡資訊失敗: " + err.Error())
				continue
			}
			if p := iip.String(); p != "" {
				availableIPs[p] = p
			}
		}
	}

	return availableIPs
}
