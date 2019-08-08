package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("curl", "https://www.google.com/", "-w", `
	{
		"size_download": %{size_download},
		"speed_download": %{speed_download},
		"time_namelookup": %{time_namelookup},
		"time_connect": %{time_connect},
		"time_appconnect": %{time_appconnect},
		"time_pretransfer": %{time_pretransfer},
		"time_redirect": %{time_redirect},
		"time_starttransfer": %{time_starttransfer},
		"time_total": %{time_total}
	}
	`,
		"-s", "-f", "-o", "/dev/null",
	)

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		log.Println("Err ==> ", stderr.String())
		log.Fatal("Run Error => ", err)
	}
	log.Println("Out ==> ", stdout.String())

	data := map[string]float64{}
	err = json.Unmarshal(stdout.Bytes(), &data)
	if err != nil {
		log.Fatal("JSON Error => ", err)
	}

	total := data["time_total"]
	nameLookup := data["time_namelookup"]
	tcpConnect := data["time_connect"] - data["time_namelookup"]
	sslConnect := data["time_appconnect"] - data["time_connect"]
	preTransfer := data["time_pretransfer"] - data["time_appconnect"]
	if data["time_appconnect"] == 0 {
		sslConnect = 0
		preTransfer = data["time_pretransfer"] - data["time_connect"]
	}
	redirect := data["time_redirect"]
	serverHandle := data["time_starttransfer"] - data["time_pretransfer"]
	returnTime := data["time_total"] - data["time_starttransfer"]

	log.Printf(`
		總時間 %f 
		-> 解析網址 %f (%.3f％)
		-> TCP握手 %f (%.3f％)
		-> SSL檢查 %f (%.3f％)
		-> 傳入資料 %f (%.3f％)
		-> 轉導 %f (%.3f％)
		-> 交給Server處理 %f (%.3f％)
		-> 回傳過程 %f (%.3f％)
	`,
		total,
		nameLookup, nameLookup/total*100,
		tcpConnect, tcpConnect/total*100,
		sslConnect, sslConnect/total*100,
		preTransfer, preTransfer/total*100,
		redirect, redirect/total*100,
		serverHandle, serverHandle/total*100,
		returnTime, returnTime/total*100,
	)
}
