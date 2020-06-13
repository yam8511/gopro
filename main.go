package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.Int("p", 2020, "監聽port")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.StaticFS("/", http.Dir("."))

	var l net.Listener
	var err error
	for {
		l, err = net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			if strings.Contains(err.Error(), "address already in use") {
				*port++
				continue
			}

			log.Fatal(err)
		}
		break
	}

	server := http.Server{Handler: r}
	go server.Serve(l)

	runtime.Gosched()
	host := fmt.Sprintf("http://127.0.0.1:%d", *port)
	fmt.Println(host)
	openbrowser(host)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	server.Shutdown(context.Background())
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println("Open Browser Error: ", err)
	}
}
