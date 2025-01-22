package interceptor

import (
	"context"
	"log"
	"time"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"
	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/resiliency"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LogUnaryClientInterceptor() grpc.UnaryClientInterceptor {

	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("[LOGGED BY UNARY CLIENT INTERCEPTOR]", req)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// Adding metadata to client request and modifying response
func BasicUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		switch request := req.(type) {
		case hello.HelloRequest:
			request.Name = "[MODIFIED BY UNARY CLIENT INTERCEPTOR - 1]" + request.Name
		}
		// adding request metadata
		ctx = metadata.AppendToOutgoingContext(ctx,
			"my-request-metadata-key-1", "my-request-metadata-value-1")
		ctx = metadata.AppendToOutgoingContext(ctx,
			"my-request-metadata-key-2", "my-request-metadata-value-2")

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			return err
		}

		//Update response received from server
		switch response := reply.(type) {
		case hello.HelloResponse:
			response.Greet = "[MODIFIED BY UNARY CLIENT INTERCEPTOR - 2]" + response.Greet
		case resiliency.ResiliencyResponse:
			response.DummyString = "[MODIFIED BY UNARY CLIENT INTERCEPTOR - 2]" + response.DummyString
		}
		return nil
	}
}

func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx, _ := context.WithTimeout(ctx, timeout)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

func LogStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Println("[LOGGED BY CLIENT STREAM INTERCEPTOR]", method)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

// Interceptor to modify the client metadata
type InterceptedClientStream struct {
	grpc.ClientStream
}

func BasicClientStreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

		ctx = metadata.AppendToOutgoingContext(ctx,
			"my-request-metadata-key-1", "my-request-metadata-value-1")
		ctx = metadata.AppendToOutgoingContext(ctx,
			"my-request-metadata-key-2", "my-request-metadata-value-2")

		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			log.Printf("Failed to start %v streaming call to %v: %v", desc.StreamName, method, err)
			return nil, err
		}

		interceptedClientStream := &InterceptedClientStream{
			ClientStream: clientStream,
		}
		return interceptedClientStream, nil
	}
}

// Modify request message
func (interceptedClientStream *InterceptedClientStream) SendMsg(msg interface{}) error {
	switch request := msg.(type) {
	case *hello.HelloRequest:
		request.Name = "[MODIFIED BY CLIENT STREAM INTERCEPTOR - 4] " + request.Name
	}
	return interceptedClientStream.ClientStream.SendMsg(msg)
}

func (interceptedClientStream *InterceptedClientStream) Recv(msg interface{}) error {
	err := interceptedClientStream.ClientStream.RecvMsg(msg)
	if err != nil {
		return err
	}
	switch response := msg.(type) {
	case *hello.HelloResponse:
		response.Greet = "[MODIFIED BY CLIENT STREAM INTERCEPTOR - 5] " + response.Greet
	case *resiliency.ResiliencyResponse:
		response.DummyString = "[MODIFIED BY CLIENT STREAM INTERCEPTOR - 6] " + response.DummyString
	}
	return nil
}

func TimeoutStreamClientIntereptor(timeout time.Duration) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		newCtx, _ := context.WithTimeout(ctx, timeout)
		return streamer(newCtx, desc, cc, method, opts...)
	}
}
