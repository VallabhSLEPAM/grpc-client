package resiliency

import (
	"context"
	"io"
	"log"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/resiliency"
	"github.com/VallabhSLEPAM/grpc-client/internal/port"
	"google.golang.org/grpc"
)

type ResiliencyAdapter struct {
	resiliencyClientPort port.ResiliencyClientPort
}

func NewResiliencyAdapter(conn *grpc.ClientConn) (*ResiliencyAdapter, error) {
	client := resiliency.NewResiliencyServiceClient(conn)
	return &ResiliencyAdapter{
		resiliencyClientPort: client,
	}, nil
}

func (adapter ResiliencyAdapter) UnaryResiliency(ctx context.Context, minDelay, maxDelay int, statusCode []uint32) (*resiliency.ResiliencyResponse, error) {
	resiliencyRequest := resiliency.ResiliencyRequest{
		MinDelaySecond: int32(minDelay),
		MaxDelaySecond: int32(maxDelay),
		StatusCodes:    statusCode,
	}

	res, err := adapter.resiliencyClientPort.UnaryResiliency(ctx, &resiliencyRequest)
	if err != nil {
		log.Println("Error on UnaryResiliency:", err)
		return nil, err
	}
	return res, nil
}

func (adapter ResiliencyAdapter) ServerResiliency(ctx context.Context, minDelay, maxDelay int, statusCode []uint32) {
	resiliencyRequest := resiliency.ResiliencyRequest{
		MinDelaySecond: int32(minDelay),
		MaxDelaySecond: int32(maxDelay),
		StatusCodes:    statusCode,
	}

	reslResp, err := adapter.resiliencyClientPort.ServerResiliency(ctx, &resiliencyRequest)
	if err != nil {
		log.Fatalln("Error on ServerResiliency:", err)
	}

	for {
		res, err := reslResp.Recv()
		if err == io.EOF {
			log.Println("Server request cancelled ServerResiliency")
			break

		}
		if err != nil {
			log.Fatalln("Error on ServerResiliency:", err)
		}
		log.Println("Response from server: ", res.DummyString)
	}
}

func (adapter ResiliencyAdapter) ClientResiliency(ctx context.Context, minDelay, maxDelay int, statusCode []uint32, count int) {

	respStream, err := adapter.resiliencyClientPort.ClientResiliency(ctx)
	if err != nil {
		log.Fatalln("Error on ClientResiliency:", err)
	}

	for i := 0; i < count; i++ {
		resiliencyRequest := resiliency.ResiliencyRequest{
			MinDelaySecond: int32(minDelay),
			MaxDelaySecond: int32(maxDelay),
			StatusCodes:    statusCode,
		}
		respStream.Send(&resiliencyRequest)
	}
	resp, err := respStream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Error on ClientResiliency:", err)
	}
	log.Println(resp.DummyString)
}

func (adapter ResiliencyAdapter) BiDirectionalResiliency(ctx context.Context, minDelay, maxDelay int, statusCode []uint32, count int) {

	respStream, err := adapter.resiliencyClientPort.BiDirectionalResiliency(ctx)
	if err != nil {
		log.Fatalln("Error on BiDirectionalResiliency:", err)
	}
	respChan := make(chan struct{})
	go func() {
		for i := 0; i < count; i++ {
			resiliencyRequest := resiliency.ResiliencyRequest{
				MinDelaySecond: int32(minDelay),
				MaxDelaySecond: int32(maxDelay),
				StatusCodes:    statusCode,
			}
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
				log.Fatalln("Error on BiDirectionalResiliency:", err)
			}
			log.Printf(res.DummyString)
		}
		close(respChan)

	}()
	<-respChan
}
