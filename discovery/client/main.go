package main

import (
	"fmt"
	"time"

	"github.com/grpc/balancer"
	"github.com/grpc/balancer/pb"

	"golang.org/x/net/context"
	"github.com/grpc"
	"github.com/grpc/resolver"
)

func main() {
	//整个的服务发现和负载均衡都交给了客户端来做，这里的底层代码很复杂，尤其是涉及到grpc部分。
	r := balancer.NewResolver("localhost:2378")
	resolver.Register(r)
//客户端使用了自实现的命名解析器etcdResolver。使用的负载均衡器为grpc自带的round robin。
	conn, err := grpc.Dial(r.Scheme()+"://author/project/test", grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewHelloServiceClient(conn)

	for {
		resp, err := client.Echo(context.Background(), &pb.Payload{Data: "hello"}, grpc.FailFast(true))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp)
		}

		<-time.After(time.Second)
	}
}