package main

import (
	"log"

	mygrpc "github.com/seminhnva/my-grpc-go-server/internal/adapter/grpc"
	app "github.com/seminhnva/my-grpc-go-server/internal/application"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(log.Writer())

	hs := &app.HelloService{}
	bs := &app.BankSerivce{}

	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)
	grpcAdapter.Run()

}
