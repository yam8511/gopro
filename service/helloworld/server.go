package helloworld

import (
	"context"
	"errors"
	"fmt"
	pb "gopro/pb/helloworld"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (*server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.GetName() == "Mr. 3" {
		return nil, errors.New("禁止進入")
	}
	msg := "Hello " + in.GetName()
	log.Println(msg)
	time.Sleep(time.Second * 5)
	res := &pb.HelloReply{Message: msg}
	log.Println("Bye~ ", in.GetName())
	return res, nil
}

func (*server) KeepHello(stream pb.Greeter_KeepHelloServer) error {
	in, err := stream.Recv()
	if err != nil {
		return err
	}

	msg := "Hello~ " + in.GetName()
	for i := 0; i < 5; i++ {
		log.Println(msg, i)
		err = stream.Send(&pb.HelloReply{Message: fmt.Sprint(msg, i)})
		if err != nil {
			return err
		}
	}
	log.Println("Finish")

	return nil
}

const (
	port = ":50052"
)

func RunServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	log.Println("監聽:", lis.Addr().String())
	<-sig
	log.Println("接收訊號，等待關閉")
	s.GracefulStop()
	log.Println("Bye~")
}
