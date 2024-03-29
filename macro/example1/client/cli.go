package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/registry/mdns"
)

// log wrapper logs every time a request is made
type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Printf("[wrapper] client request service: %s method: %s\n", req.Service(), req.Method())
	return l.Client.Call(ctx, req, rsp)
}

// Implements client.Wrapper as logWrapper
func logWrap(c client.Client) client.Client {
	return &logWrapper{c}
}

func main() {
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.Registry(mdns.NewRegistry()),
		// wrap the client
		micro.WrapClient(logWrap),
	)

	service.Init()

	greeter := proto.NewGreeterService("greeter", service.Client())

	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "John"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Greeting)

}