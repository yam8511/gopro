package helloworld

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	pb "gopro/pb/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":50052"
)

func RunClient() {
	conn, err := grpc.Dial(os.Getenv("IP")+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		log.Println("Init", i)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Println("Start", i)
			res, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: fmt.Sprint("Mr. ", i)})
			log.Println("End", i, err)

			s, _ := status.FromError(err)
			switch s.Code() {
			case codes.OK:
				log.Println("OK: ", res.GetMessage())
			case codes.Unavailable:
				log.Println("Unavailable", i, ": 服務連不到: ", s.Message())
			case codes.Unknown:
				log.Println("Unknown", i, ": ", s.Message())
			default:
				log.Println("捕獲未定義錯誤", i, ": ", s.Message())
			}
		}(i)
	}

	wg.Wait()
}
