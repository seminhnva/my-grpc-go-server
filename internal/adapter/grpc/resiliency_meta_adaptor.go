package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/seminhnva/my-grpc-proto/protogen/go/resiliency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func dummyRequestMetadata(ctx context.Context) {
	if requestMetadata, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println("Request data")
		for k, v := range requestMetadata {
			log.Printf(" %v : %v\n", k, v)
		}
	} else {
		log.Println("Request metadata not found")
	}
}

func dummyResponseMetadata() metadata.MD {
	md := map[string]string{
		"grpc-server-time":     fmt.Sprint(time.Now().Format("15:04:05")),
		"grpc-server-location": "vietname",
		"grpc-server-uuid":     uuid.New().String(),
	}
	return metadata.New(md)
}

func (a *GrpcAdapter) UnaryResiliencyWithMetadata(ctx context.Context, req *resiliency.ResiliencyRequest) (
	*resiliency.ResiliencyResponse, error) {
	log.Println("UnaryResiliencyWithMetadata called")
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

	dummyRequestMetadata(ctx)

	if errStatus := generateErrStatus(sts); errStatus != nil {
		return nil, errStatus
	}

	grpc.SendHeader(ctx, dummyResponseMetadata())
	return &resiliency.ResiliencyResponse{
		DummyString: str,
	}, nil

}
func (a *GrpcAdapter) ServerStreamResiliencyWithMetadata(req *resiliency.ResiliencyRequest,
	stream grpc.ServerStreamingServer[resiliency.ResiliencyResponse]) error {
	log.Println("ServerStreamResiliencyWithMetadat called")

	ctx := stream.Context()
	dummyRequestMetadata(ctx)
	if err := stream.SendHeader(dummyResponseMetadata()); err != nil {
		log.Println("Error while sending respone metadata: ", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("ServerStreamResiliencyWithMetadat stopped: %v", ctx.Err())
			return nil
		default:
			str, sts := a.resililencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.StatusCodes)
			if ctx.Err() != nil {
				log.Printf("ServerStreamResiliencyWithMetadat stopped: %v", ctx.Err())
				return nil
			}
			if errStatus := generateErrStatus(sts); errStatus != nil {
				return errStatus
			}
			err := stream.Send(&resiliency.ResiliencyResponse{
				DummyString: str,
			})
			if err != nil {
				if ctx.Err() != nil {
					log.Printf("ServerStreamResiliencyWithMetadat stopped: %v", ctx.Err())
					return nil
				}
				return err
			}
		}

	}

}
func (a *GrpcAdapter) ClientStreamResiliencyWithMetadata(stream grpc.ClientStreamingServer[resiliency.ResiliencyRequest,
	resiliency.ResiliencyResponse]) error {
	log.Println("ClientStreamResiliencyWithMetadat called")
	i := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := resiliency.ResiliencyResponse{
				DummyString: fmt.Sprintf("Received %v request from client", strconv.Itoa(i)),
			}

			if err := stream.SendHeader(dummyResponseMetadata()); err != nil {
				log.Println("Error while sending respone metadata: ", err)
			}

			return stream.SendAndClose(&res)
		}
		context := stream.Context()
		dummyRequestMetadata(context)

		if req != nil {
			_, sts := a.resililencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.GetStatusCodes())
			if errStatus := generateErrStatus(sts); errStatus != nil {
				return errStatus
			}
		}
		i = i + 1
	}

}
func (a *GrpcAdapter) BiDirectionResiliencyWithMetadata(stream grpc.BidiStreamingServer[resiliency.ResiliencyRequest,
	resiliency.ResiliencyResponse]) error {
	log.Println("BiDirectionResiliencyWithMetadat called")

	context := stream.Context()
	if err := stream.SendHeader(dummyResponseMetadata()); err != nil {
		log.Println("Error while sending respone metadata: ", err)
	}

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
			dummyRequestMetadata(context)
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
