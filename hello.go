package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cin = make(chan string)
var cout = make(chan string)

// Arith 數學運算
type Arith int

// Sum 總和
func (arith *Arith) Sum(a *int, b *int) (err error) {
	time.Sleep(time.Second * 10)
	*b = *a + 1
	return
}

func main() {
	fmt.Println("Hello World")

	jsonrpcAddress := flag.String("jsonrpc_address", ":8000", "the address for server listening")
	isClient := flag.Bool("c", false, "if run client")
	flag.Parse()

	if *isClient {
		runJSONRPCClient(*jsonrpcAddress)
		return
	}

	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	online := 0
	done := false

	go func() {
		rpc.RegisterName("Arith", new(Arith))
		for {
			conn, err := l.Accept()
			if err != nil {
				if done {
					return
				}
				log.Println("Error: accept jsonrpc connection ->", err)
				continue
			}
			log.Println("Accept jsonrpc connection")
			// 設定連線timeout
			// conn.SetDeadline(time.Now().Add(time.Second * time.Duration(*timeout)))
			go func(conn net.Conn) {
				ip := conn.RemoteAddr().String()
				cin <- ip
				jsonrpc.ServeConn(conn)
				cout <- ip
			}(conn)
		}
	}()

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

func runJSONRPCClient(address string) {
	client, err := jsonrpc.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var args interface{}
	args = 1
	var sum int
	err = client.Call("Arith.Sum", args, &sum)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: req -> %v , res -> %v\n", args, sum)
}
