package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/seminhnva/my-grpc-proto/protogen/go/resiliency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type resiliencyServiceFunc func(int32, int32, []uint32) (string, uint32)

func (f resiliencyServiceFunc) GenerateResiliency(minDelaySecond, maxDelaySecond int32, statusCodes []uint32) (string, uint32) {
	return f(minDelaySecond, maxDelaySecond, statusCodes)
}

type serverStreamingTestStream struct {
	ctx       context.Context
	sendErr   error
	sendCalls int
}

func (s *serverStreamingTestStream) Send(*resiliency.ResiliencyResponse) error {
	s.sendCalls++
	return s.sendErr
}

func (s *serverStreamingTestStream) SetHeader(metadata.MD) error  { return nil }
func (s *serverStreamingTestStream) SendHeader(metadata.MD) error { return nil }
func (s *serverStreamingTestStream) SetTrailer(metadata.MD)       {}
func (s *serverStreamingTestStream) Context() context.Context     { return s.ctx }
func (s *serverStreamingTestStream) SendMsg(any) error            { return nil }
func (s *serverStreamingTestStream) RecvMsg(any) error            { return nil }

func TestGenerateErrStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode uint32
		wantCode   codes.Code
		wantErr    bool
	}{
		{name: "OK", statusCode: 0, wantCode: codes.OK},
		{name: "Unavailable", statusCode: 14, wantCode: codes.Unavailable, wantErr: true},
		{name: "unknown numeric status", statusCode: 99, wantCode: codes.Internal, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generateErrStatus(tt.statusCode)
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("generateErrStatus(%d) = %v, want nil", tt.statusCode, err)
				}
				return
			}

			if err == nil {
				t.Fatalf("generateErrStatus(%d) = nil, want error", tt.statusCode)
			}
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("status.Code(err) = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func TestServerStreamResiliencyStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stream := &serverStreamingTestStream{ctx: ctx}
	service := resiliencyServiceFunc(func(int32, int32, []uint32) (string, uint32) {
		cancel()
		return "generated response", 0
	})
	adapter := &GrpcAdapter{resililencyService: service}

	err := adapter.ServerStreamResiliency(&resiliency.ResiliencyRequest{
		MinDelaySecond: 1,
		MaxDelaySecond: 2,
		StatusCodes:    []uint32{0},
	}, stream)
	if err != nil {
		t.Fatalf("ServerStreamResiliency() error = %v, want nil after cancellation", err)
	}
	if stream.sendCalls != 0 {
		t.Errorf("Send called %d times after cancellation, want 0", stream.sendCalls)
	}
}

func TestServerStreamResiliencyReturnsSendError(t *testing.T) {
	sendErr := errors.New("send failed")
	stream := &serverStreamingTestStream{
		ctx:     context.Background(),
		sendErr: sendErr,
	}
	service := resiliencyServiceFunc(func(int32, int32, []uint32) (string, uint32) {
		return "generated response", 0
	})
	adapter := &GrpcAdapter{resililencyService: service}

	err := adapter.ServerStreamResiliency(&resiliency.ResiliencyRequest{
		MinDelaySecond: 1,
		MaxDelaySecond: 2,
		StatusCodes:    []uint32{0},
	}, stream)
	if !errors.Is(err, sendErr) {
		t.Errorf("ServerStreamResiliency() error = %v, want %v", err, sendErr)
	}
}
