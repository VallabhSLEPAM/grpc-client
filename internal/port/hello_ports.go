package port

import (
	"context"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"
	"google.golang.org/grpc"
)

type HelloClientPort interface {
	SayHello(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error)
}
