package port

import (
	"context"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/resiliency"
	"google.golang.org/grpc"
)

type ResiliencyClientPort interface {
	UnaryResiliency(ctx context.Context, in *resiliency.ResiliencyRequest, opts ...grpc.CallOption) (*resiliency.ResiliencyResponse, error)
	ServerResiliency(ctx context.Context, in *resiliency.ResiliencyRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[resiliency.ResiliencyResponse], error)
	ClientResiliency(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[resiliency.ResiliencyRequest, resiliency.ResiliencyResponse], error)
	BiDirectionalResiliency(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[resiliency.ResiliencyRequest, resiliency.ResiliencyResponse], error)
}

type ResiliencyMetadataClientPort interface {
	UnaryResiliencyWithMetadata(ctx context.Context, in *resiliency.ResiliencyRequest, opts ...grpc.CallOption) (*resiliency.ResiliencyResponse, error)
	ServerResiliencyWithMetadata(ctx context.Context, in *resiliency.ResiliencyRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[resiliency.ResiliencyResponse], error)
	ClientResiliencyWithMetadata(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[resiliency.ResiliencyRequest, resiliency.ResiliencyResponse], error)
	BiDirectionalResiliencyWithMetadata(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[resiliency.ResiliencyRequest, resiliency.ResiliencyResponse], error)
}
