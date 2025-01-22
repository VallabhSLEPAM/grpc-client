package resiliency

import (
	"context"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/resiliency"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func sampleRequestMetadata() metadata.MD {
	md := map[string]string{
		"grpc-client-time":  time.Now().Format("20:40:00"),
		"grpc-client-os":    runtime.GOOS,
		"grpc-request-uuid": uuid.New().String(),
	}
	return metadata.New(md)
}

func sampleResponseMetadata(md metadata.MD) {
	if md.Len() == 0 {
		fmt.Println("No metadata found.")
	} else {
		fmt.Println("Response metadata:")

		for k, v := range md {
			fmt.Printf("Key: %v,Value:%v\n", k, v)
		}
	}

}

func (adapter ResiliencyAdapter) UnaryResiliencyWithMetadata(ctx context.Context, minDelay, maxDelay int, statusCode []uint32) (*resiliency.ResiliencyResponse, error) {
	resiliencyRequest := resiliency.ResiliencyRequest{
		MinDelaySecond: int32(minDelay),
		MaxDelaySecond: int32(maxDelay),
		StatusCodes:    statusCode,
	}
	ctx = metadata.NewOutgoingContext(ctx, sampleRequestMetadata())

	var responseMetadata metadata.MD
	res, err := adapter.resiliencyMetadataClientPort.UnaryResiliencyWithMetadata(ctx, &resiliencyRequest, grpc.Header(&responseMetadata))
	if err != nil {
		log.Println("Error on UnaryResiliencyWithMetadata:", err)
		return nil, err
	}
	sampleResponseMetadata(responseMetadata)
	return res, nil
}

func (adapter ResiliencyAdapter) ServerResiliencyWithMetadata(ctx context.Context, minDelay, maxDelay int, statusCode []uint32) {
	resiliencyRequest := resiliency.ResiliencyRequest{
		MinDelaySecond: int32(minDelay),
		MaxDelaySecond: int32(maxDelay),
		StatusCodes:    statusCode,
	}
	ctx = metadata.NewOutgoingContext(ctx, sampleRequestMetadata())

	reslResp, err := adapter.resiliencyClientPort.ServerResiliency(ctx, &resiliencyRequest)
	if err != nil {
		log.Fatalln("Error on ServerResiliencyWithMetadata:", err)
	}

	if respMetadata, err := reslResp.Header(); err == nil {
		sampleResponseMetadata(respMetadata)
	}

	for {
		res, err := reslResp.Recv()
		if err == io.EOF {
			log.Println("Server request cancelled ServerResiliencyWithMetadata")
			break

		}
		if err != nil {
			log.Fatalln("Error on ServerResiliencyWithMetadata:", err)
		}
		log.Println("Response from server: ", res.DummyString)
	}
}

func (adapter ResiliencyAdapter) ClientResiliencyWithMetadata(ctx context.Context, minDelay, maxDelay int, statusCode []uint32, count int) {

	respStream, err := adapter.resiliencyClientPort.ClientResiliency(ctx)
	if err != nil {
		log.Fatalln("Error on ClientResiliencyWithMetadata:", err)
	}

	for i := 0; i < count; i++ {

		ctx = metadata.NewOutgoingContext(ctx, sampleRequestMetadata())
		resiliencyRequest := resiliency.ResiliencyRequest{
			MinDelaySecond: int32(minDelay),
			MaxDelaySecond: int32(maxDelay),
			StatusCodes:    statusCode,
		}
		respStream.Send(&resiliencyRequest)
	}

	if respMetadata, err := respStream.Header(); err == nil {
		sampleResponseMetadata(respMetadata)
	}

	resp, err := respStream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Error on ClientResiliencyWithMetadata:", err)
	}
	log.Println(resp.DummyString)
}

func (adapter ResiliencyAdapter) BiDirectionalResiliencyWithMetadata(ctx context.Context, minDelay, maxDelay int, statusCode []uint32, count int) {

	respStream, err := adapter.resiliencyClientPort.BiDirectionalResiliency(ctx)
	if err != nil {
		log.Fatalln("Error on BiDirectionalResiliencyWithMetadata:", err)
	}

	if respMetadata, err := respStream.Header(); err == nil {
		sampleResponseMetadata(respMetadata)
	}

	respChan := make(chan struct{})
	go func() {
		for i := 0; i < count; i++ {
			resiliencyRequest := resiliency.ResiliencyRequest{
				MinDelaySecond: int32(minDelay),
				MaxDelaySecond: int32(maxDelay),
				StatusCodes:    statusCode,
			}
			ctx = metadata.NewOutgoingContext(ctx, sampleRequestMetadata())
			respStream.Send(&resiliencyRequest)
		}
		respStream.CloseSend()
	}()

	go func() {
		for {
			res, err := respStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln("Error on BiDirectionalResiliencyWithMetadata:", err)
			}
			log.Printf(res.DummyString)
		}
		close(respChan)

	}()
	<-respChan
}
