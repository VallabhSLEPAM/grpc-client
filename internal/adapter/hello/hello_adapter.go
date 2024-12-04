package adapter

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"
	"github.com/VallabhSLEPAM/grpc-client/internal/port"
	"google.golang.org/grpc"
)

type HelloAdapter struct {
	helloClient port.HelloClientPort
}

// Just a wrapper to create service client and return it
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

func (a *HelloAdapter) SayHelloServerStream(ctx context.Context, name string) {
	helloRequest := &hello.HelloRequest{
		Name: "Vallabh",
	}

	// Making the server RPC call
	greetStream, err := a.helloClient.HelloServerStream(ctx, helloRequest)
	if err != nil {
		log.Fatalln("Error on SayHelloServerStream: ", err)
	}

	for {
		greet, err := greetStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error on HelloServerStream while receiving: ", err)
		}
		log.Println(greet.Greet)

	}
}

func (a *HelloAdapter) SayHelloClientStream(ctx context.Context, names []string) {

	greetStream, err := a.helloClient.HelloClientStream(ctx)
	if err != nil {
		log.Fatalln("Error on SayHelloClientStream: ", err)
	}

	for _, name := range names {
		req := &hello.HelloRequest{
			Name: name,
		}
		err = greetStream.Send(req)
		if err != nil {
			log.Fatalf("Error sending request with message: %v. err: %v\n", req, err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	resp, err := greetStream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Error getting response in SayHelloClientStream: ", err)
	}
	log.Println(resp.Greet)

}

func (a *HelloAdapter) SayHelloContinuous(ctx context.Context, names []string) {

	stream, err := a.helloClient.HelloContinuous(ctx)
	if err != nil {
		log.Fatalf("Error getting client stream: %v\n", err)
	}

	greetChan := make(chan struct{})
	go func() {

		for _, name := range names {
			stream.Send(&hello.HelloRequest{
				Name: name,
			})
		}
		stream.CloseSend()
	}()

	go func() {

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error getting response from server stream: %v\n", err)
			}
			log.Println("Received response from server stream:", resp.Greet)

		}
		close(greetChan)
	}()

	<-greetChan

}
