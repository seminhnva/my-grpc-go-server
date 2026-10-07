package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/seminhnva/my-grpc-go-server/internal/port"
	"github.com/seminhnva/my-grpc-proto/protogen/go/bank"
	"github.com/seminhnva/my-grpc-proto/protogen/go/hello"
	"google.golang.org/grpc"
)

type GrpcAdapter struct {
	helloService port.HelloSerivcePort
	bankService  port.BankServicePort
	grpcPort     int
	server       *grpc.Server
	hello.HelloServiceServer
	bank.BankServiceServer
}

func NewGrpcAdapter(helloService port.HelloSerivcePort, bankService port.BankServicePort, grpcPort int) *GrpcAdapter {
	return &GrpcAdapter{
		helloService: helloService,
		bankService:  bankService,
		grpcPort:     grpcPort,
	}
}

func (a *GrpcAdapter) Run() {
	var err error
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d", a.grpcPort)
	}

	grpcServer := grpc.NewServer()
	a.server = grpcServer

	hello.RegisterHelloServiceServer(grpcServer, a)
	bank.RegisterBankServiceServer(grpcServer, a)
	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve grpc on port %d : %v\n", a.grpcPort, err)

	}
}

func (a *GrpcAdapter) Stop() {
	a.server.Stop()
}
