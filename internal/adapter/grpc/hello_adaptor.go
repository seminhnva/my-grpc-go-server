package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
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

func (a *GrpcAdapter) SayHelloToEveryone(stream grpc.ClientStreamingServer[hello.HelloRequest, hello.HelloResponse]) error {
	res := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&hello.HelloResponse{
				Greet: res,
			})
		}
		greet := a.helloService.GenerateGreet(req.Name)
		res += greet + " "
	}
}

func (a *GrpcAdapter) SayHelloContinuous(stream grpc.BidiStreamingServer[hello.HelloRequest, hello.HelloResponse]) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		greet := a.helloService.GenerateGreet(req.Name)
		if err = stream.Send(&hello.HelloResponse{Greet: greet}); err != nil {
			log.Fatalln("Error while say HelloContinious", err)
		}
	}
}
