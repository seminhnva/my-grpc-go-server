package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/seminhnva/my-grpc-proto/protogen/go/hello"
	"google.golang.org/grpc"
)

func (a *GrpcAdapter) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	greet := a.helloService.GenerateGreet(req.Name)

	return &hello.HelloResponse{
		Greet: greet,
	}, nil
}

func (a *GrpcAdapter) SayManyHello(req *hello.HelloRequest, stream grpc.ServerStreamingServer[hello.HelloResponse]) error {
	for i := 0; i <= 10; i++ {
		greet := a.helloService.GenerateGreet(req.Name)
		res := fmt.Sprintf("[%d] %s", i, greet)

		stream.Send(
			&hello.HelloResponse{
				Greet: res,
			},
		)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
