package adapter

import (
	"context"
	"log"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"
	"github.com/VallabhSLEPAM/grpc-client/internal/port"
	"google.golang.org/grpc"
)

type HelloAdapter struct {
	helloClient port.HelloClientPort
}

func NewHelloAdapter(conn *grpc.ClientConn) (*HelloAdapter, error) {
	client := hello.NewHelloServiceClient(conn)
	return &HelloAdapter{
		helloClient: client,
	}, nil
}

func (a *HelloAdapter) SayHello(ctx context.Context, name string) (*hello.HelloResponse, error) {
	helloRequest := &hello.HelloRequest{
		Name: "Vallabh",
	}

	greet, err := a.helloClient.SayHello(ctx, helloRequest)
	if err != nil {
		log.Fatal("Error on SayHello: ", err)
	}
	return greet, nil

}
