package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"sync"
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
	time.Sleep(time.Second * 1)
	*sum = args.A + args.B
	return nil
}

// Diff 差和
func (t *Arith) Diff(args *Args, diff *int) error {
	*diff = args.A - args.B
	return nil
}

func main() {
	rpcAddress := flag.String("rpc_address", ":50051", "the address for server listening")
	jsonrpcAddress := flag.String("jsonrpc_address", ":50052", "the address for server listening")
	httpAddress := flag.String("http_address", ":8888", "the address for server listening")
	timeout := flag.Int("timeout", 5, "Connection Timeout")
	isClient := flag.Bool("c", false, "if run client")
	flag.Parse()

	if *isClient {
		runRPCClient(*rpcAddress)
		runJSONRPCClient(*jsonrpcAddress)
		return
	}

	arith := new(Arith)
	// RPC
	rpc.RegisterName("arith", arith)
	l, e := net.Listen("tcp", *rpcAddress)
	if e != nil {
		log.Fatal("rpc listen error:", e)
	}

	// JSON-RPC
	j, e := net.Listen("tcp", *jsonrpcAddress)
	if e != nil {
		log.Fatal("jsonrpc listen error:", e)
	}

	// HTTP
	h, e := net.Listen("tcp", *httpAddress)
	if e != nil {
		log.Fatal("http listen error:", e)
	}
	httpServer := http.Server{
		Handler: new(Handler),
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
			j.Close()
			httpServer.Close()
		}
	}(l)

	wg := new(sync.WaitGroup)
	wg.Add(3)

	// RPC
	go func() {
		defer wg.Done()
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
			// 設定連線timeout
			conn.SetDeadline(time.Now().Add(time.Second * time.Duration(*timeout)))
			go rpc.ServeConn(conn)
		}
	}()

	// JSON-RPC
	go func() {
		defer wg.Done()
		for {
			conn, err := j.Accept()
			if err != nil {
				select {
				case <-c:
					return
				default:
					log.Println("Error: accept jsonrpc connection ->", err)
				}
				continue
			}
			log.Println("Accep jsonrpc connection")
			// 設定連線timeout
			conn.SetDeadline(time.Now().Add(time.Second * time.Duration(*timeout)))
			go jsonrpc.ServeConn(conn)
		}
	}()

	// HTTP
	go func() {
		defer wg.Done()
		err := httpServer.Serve(h)
		if err != nil {
			select {
			case <-c:
				return
			default:
				log.Println("Error: accept http connection ->", err)
			}
		}
	}()

	log.Println("RPC Server Listening ... ", l.Addr().Network(), l.Addr().String())
	log.Println("JSON-RPC Server Listening ... ", j.Addr().Network(), j.Addr().String())
	log.Println("HTTP Server Listening ... ", h.Addr().Network(), h.Addr().String())

	wg.Wait()
}

func runRPCClient(address string) {
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
	fmt.Printf("Arith: req -> %v , res -> %v\n", args, sum)
}

func runJSONRPCClient(address string) {
	client, err := jsonrpc.Dial("tcp", address)
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
	fmt.Printf("Arith: req -> %v , res -> %v\n", args, sum)
}

func transferJSONRPCClient(address, method string, params interface{}) (res interface{}, err error) {
	client, dialErr := jsonrpc.Dial("tcp", address)
	if dialErr != nil {
		err = dialErr
		return
	}
	defer client.Close()
	err = client.Call(method, params, &res)
	return
}
