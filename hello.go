package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cin = make(chan string)
var cout = make(chan string)

type server int

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cin <- r.RemoteAddr
	time.Sleep(time.Second * 10)
	w.Write([]byte(fmt.Sprintf("You are from %s", r.RemoteAddr)))
	cout <- r.RemoteAddr
}

func main() {
	fmt.Println("Hello World")
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		err = http.Serve(l, new(server))
		log.Println("http err ->", err)
	}()

	online := 0
	done := false
	for {
		done = delectSignal(done, sig)
		if done {
			log.Println("WAIT EXIT ...")
			for {
				log.Println("Online -> ", online)
				if online <= 0 {
					break
				}
				online, done = waitConnection(online, done, sig)
			}
			l.Close()
			log.Println("EXIT")
			return
		}

		online, done = waitConnection(online, done, sig)
	}
}

func delectSignal(done bool, sig chan os.Signal) bool {
	if done {
		return true
	}
	select {
	case <-sig:
		log.Println("Getcha")
		return true
	default:
		log.Println("Nothing")
		return false
	}
}

func waitConnection(online int, done bool, sig chan os.Signal) (int, bool) {
	log.Println("Done -> ", done)
	if done {
		select {
		case ip := <-cin:
			online++
			log.Println("CONNECT ->", ip)
		case ip := <-cout:
			online--
			log.Println("DISCONNECT ->", ip)
		}
	} else {
		select {
		case ip := <-cin:
			online++
			log.Println("CONNECT ->", ip)
		case ip := <-cout:
			online--
			log.Println("DISCONNECT ->", ip)
		case <-sig:
			log.Println("wait Getcha", done)
			return online, true
		}
	}
	return online, done
}
