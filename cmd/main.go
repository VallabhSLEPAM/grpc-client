package main

import (
	"context"
	"log"

	adapter "github.com/VallabhSLEPAM/grpc-client/internal/adapter/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	grpcClient, err := grpc.NewClient("localhost:9090", opts...)
	if err != nil {
		log.Fatalf("Cannot connect to gRPC server :%v", err)
	}

	defer grpcClient.Close()

	helloAdapter, err := adapter.NewHelloAdapter(grpcClient)
	if err != nil {
		log.Fatalf("Error while creating HelloAdapter :%v", err)
	}

	runSayHello(*helloAdapter, "Bruce Banner")
}

func runSayHello(adapter adapter.HelloAdapter, name string) {
	greet, err := adapter.SayHello(context.Background(), name)
	if err != nil {
		log.Fatalf("Error while creating HelloAdapter :%v", err)
	}
	log.Println(greet.Greet)
}
