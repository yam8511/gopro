package main

import (
	"flag"
	client_hello "gopro/client/helloworld"
	server_hello "gopro/service/helloworld"
)

func main() {
	rc := flag.Bool("c", false, "run client")
	flag.Parse()

	if *rc {
		client_hello.RunClient()
	} else {
		server_hello.RunServer()
	}
}
