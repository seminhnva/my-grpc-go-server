package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	dresil "github.com/seminhnva/my-grpc-go-server/internal/application/domain/resiliency"
	"github.com/seminhnva/my-grpc-proto/protogen/go/resiliency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func generateErrStatus(statusCode uint32) error {
	str := "Generate by server"

	if sc, ok := dresil.StatusCodeMap[statusCode]; ok && sc == codes.OK {
		return nil
	}
	if sc, ok := dresil.StatusCodeMap[statusCode]; ok {
		return status.New(sc, str).Err()
	} else {
		return status.New(codes.Internal, str).Err()
	}
}

func (a *GrpcAdapter) UnaryResiliency(ctx context.Context, req *resiliency.ResiliencyRequest) (
	*resiliency.ResiliencyResponse, error) {
	log.Println("UnaryResililency called")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "resiliency request is required")
	}
	if req.MinDelaySecond < 0 || req.MaxDelaySecond <= req.MinDelaySecond {
		return nil, status.Error(codes.InvalidArgument, "max_delay_second must be greater than min_delay_second, and min_delay_second cannot be negative")
	}
	if len(req.StatusCodes) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one status code is required")
	}

	str, sts := a.resililencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.StatusCodes)
	if errStatus := generateErrStatus(sts); errStatus != nil {
		return nil, errStatus
	}

	return &resiliency.ResiliencyResponse{
		DummyString: str,
	}, nil

}
func (a *GrpcAdapter) ServerStreamResiliency(req *resiliency.ResiliencyRequest,
	stream grpc.ServerStreamingServer[resiliency.ResiliencyResponse]) error {
	log.Println("ServerStreamResiliency called")

	context := stream.Context()
	for {
		select {
		case <-context.Done():
			log.Println("Client cancelled request")
			return nil
		default:
			str, sts := a.resililencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.StatusCodes)
			if errStatus := generateErrStatus(sts); errStatus != nil {
				return errStatus
			}
			err := stream.Send(&resiliency.ResiliencyResponse{
				DummyString: str,
			})
			if err != nil {
				log.Fatalln("Error while seding to client:", err)
			}
		}

	}

}
func (a *GrpcAdapter) ClientStreamResiliency(stream grpc.ClientStreamingServer[resiliency.ResiliencyRequest,
	resiliency.ResiliencyResponse]) error {
	log.Println("ClientStreamResiliency called")
	i := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := resiliency.ResiliencyResponse{
				DummyString: fmt.Sprintf("Received %v request from client", strconv.Itoa(i)),
			}
			return stream.SendAndClose(&res)
		}
		if req != nil {
			_, sts := a.resililencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.GetStatusCodes())
			if errStatus := generateErrStatus(sts); errStatus != nil {
				return errStatus
			}
		}
		i = i + 1
	}

}
func (a *GrpcAdapter) BiDirectionResiliency(stream grpc.BidiStreamingServer[resiliency.ResiliencyRequest,
	resiliency.ResiliencyResponse]) error {
	log.Println("BiDirectionResiliency called")

	context := stream.Context()
	for {
		select {
		case <-context.Done():
			log.Println("Client cancelled request")
			return nil
		default:
			req, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				log.Fatalln("Error while reading from client:", err)
			}
			str, sts := a.resililencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.StatusCodes)
			if errStatus := generateErrStatus(sts); errStatus != nil {
				return errStatus
			}
			err = stream.Send(&resiliency.ResiliencyResponse{
				DummyString: str,
			})
			if err != nil {
				log.Fatalln("Error while seding to client:", err)
			}
		}

	}

}
