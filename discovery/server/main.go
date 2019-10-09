package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"grpc/balancer"
	"../pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const svcName = "project/test"

var addr = "127.0.0.1:50051"

func main() {
	//服务端的原理很简单，启动的时候调用注册的方法，停止的时候注销
	flag.StringVar(&addr, "addr", addr, "addr to lis")
	flag.Parse()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	defer s.GracefulStop()

	pb.RegisterHelloServiceServer(s, &hello{})

	go balancer.Register("127.0.0.1:2379", svcName, addr, 5)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		balancer.UnRegister(svcName, addr)

		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

type hello struct {
}

func (*hello) Echo(ctx context.Context, req *pb.Payload) (*pb.Payload, error) {
	req.Data = req.Data + ", from:" + addr
	return req, nil
}