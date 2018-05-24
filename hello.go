package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

// Arith 數學運算
type Arith int

// Args 參數
type Args struct {
	A, B int
}

// Sum 總和
func (t *Arith) Sum(args *Args, sum *int) error {
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
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", *address)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	go func(l net.Listener) {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		select {
		case s := <-sig:
			log.Printf("\nReceived a signal, shutdown by ... %v", s)
		}
		l.Close()
	}(l)

	log.Println("Server Listening ... ", *address)
	err := http.Serve(l, nil)
	if err != nil {
		log.Println(err)
	}

}

func runClient(address string) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	args := &Args{7, 8}
	var sum int
	err = client.Call("Arith.Sum", args, &sum)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d + %d = %d", args.A, args.B, sum)
}
