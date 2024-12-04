package port

import (
	"context"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"
	"google.golang.org/grpc"
)

type HelloClientPort interface {
	SayHello(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error)
	HelloServerStream(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[hello.HelloResponse], error)
	HelloClientStream(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[hello.HelloRequest, hello.HelloResponse], error)
	HelloContinuous(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[hello.HelloRequest, hello.HelloResponse], error)
}
