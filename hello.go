package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Arith 數學運算
type Arith int

// Args 參數
type Args struct {
	A, B int
}

// Sum 總和
func (t *Arith) Sum(args *Args, sum *int) error {
	time.Sleep(time.Second * 3)
	*sum = args.A + args.B
	return nil
}

// Diff 差和
func (t *Arith) Diff(args *Args, diff *int) error {
	*diff = args.A - args.B
	return nil
}

func main() {
	address := flag.String("address", ":50051", "the address for server listening")
	advertise := flag.String("advertise", "", "the address for call")
	isClient := flag.Bool("c", false, "if run client")
	flag.Parse()

	if *isClient {
		runClient(*address)
		return
	}

	if *advertise == "" {
		advertise = address
	}

	arith := new(Arith)
	rpc.RegisterName("arith", arith)
	l, e := net.Listen("tcp", *address)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	c := make(chan int)
	go func(l net.Listener) {
		select {
		case s := <-sig:
			log.Printf("... Receive signal, shutdown by ... %v", s)
			close(c)
			l.Close()
		}
	}(l)

	log.Println("Server Listening ... ", *address)

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-c:
				return
			default:
				log.Println("Error: accept rpc connection ->", err)
			}
			continue
		}
		log.Println("Accep rpc connection")
		rpc.ServeConn(conn)
	}

}

func runClient(address string) {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var args interface{}
	args = &Args{7, 8}
	var sum int
	err = client.Call("arith.Sum", args, &sum)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: req -> %v , res -> %v", args, sum)
}
