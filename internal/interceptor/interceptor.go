package interceptor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/seminhnva/my-grpc-proto/protogen/go/hello"
	"github.com/seminhnva/my-grpc-proto/protogen/go/resiliency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func protoMap(m interface{}) (map[string]interface{}, error) {
	pm, ok := m.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("not a proto message: %T", m)
	}
	b, err := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}.Marshal(pm)
	if err != nil {
		return nil, err
	}

	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		dataReq, err := protoMap(req)
		if err != nil {
			return nil, err
		}
		log.Printf("[UNARY-LOG-INTERCEPTER]:Inconming Request: %v \n", dataReq)
		ModifyRequest(req)
		ctx = modifyIncoming(ctx)
		res, err := handler(ctx, req)
		ModifyResponse(res)
		return res, err
	}
}

func ModifyRequest(req interface{}) {
	switch rq := req.(type) {
	case *hello.HelloRequest:
		rq.Name = "[INTERCEPTOR] " + rq.Name
	}

}
func ModifyResponse(res interface{}) {
	switch response := res.(type) {
	case *hello.HelloResponse:
		response.Greet = "[MODIFIED]" + response.Greet
	case *resiliency.ResiliencyResponse:
		response.DummyString = "[MODIFIED]" + response.DummyString
	}
}
func modifyIncoming(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md.Append("key-to-add", "value-to-add") // thêm
	md.Set("key-to-modify", "new-value")    // ghi đè
	md.Delete("key-to-delete")              // xóa key user gửi\

	_ = grpc.SetHeader(ctx, md)
	return metadata.NewIncomingContext(ctx, md)
}

type InterceptedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := modifyIncoming(stream.Context())
		serverStream := &InterceptedServerStream{ServerStream: stream, ctx: ctx}
		return handler(srv, serverStream)
	}
}
func (s *InterceptedServerStream) Context() context.Context {
	return s.ctx
}
func (s *InterceptedServerStream) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	if dataReq, err := protoMap(m); err == nil {
		log.Printf("[STREAM-LOG-INTERCEPTOR]: Incoming Request: %v \n", dataReq)
	}
	ModifyRequest(m)

	return nil
}

func (s *InterceptedServerStream) SendMsg(msg interface{}) error {
	switch response := msg.(type) {
	case *hello.HelloResponse:
		response.Greet = "[Modified by server interceptor]" + response.Greet
	case *resiliency.ResiliencyResponse:
		response.DummyString = "[Modified by server interceptor]" + response.DummyString
	}
	return s.ServerStream.SendMsg(msg)
}
