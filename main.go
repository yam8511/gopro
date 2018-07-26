package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	path := os.Getenv("GITEA_CUSTOM")
	ini := path + "/conf/app.ini"
	content, err := ioutil.ReadFile(ini)
	if err != nil {
		log.Fatalln("Read Error --->", err)
		return
	}

	port := "3000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	content = []byte(strings.Replace(string(content), "???", port, -1))
	err = ioutil.WriteFile(ini, content, 0777)
	if err != nil {
		log.Fatalln("Write Error --->", err)
		return
	}

	cmd := exec.Command("./gitea", "web")
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		log.Println("Run Error --->", err)
		return
	}
}
