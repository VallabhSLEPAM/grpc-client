package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	bankadapter "github.com/VallabhSLEPAM/grpc-client/internal/adapter/bank"
	adapter "github.com/VallabhSLEPAM/grpc-client/internal/adapter/hello"
	"github.com/VallabhSLEPAM/grpc-client/internal/application/domain/bank"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Create a gRPC client with TLS credentials
	grpcClient, err := grpc.NewClient("localhost:9090", opts...)
	if err != nil {
		log.Fatalf("Cannot connect to gRPC server :%v", err)
	}

	//defer to close the gRPC client connection
	defer grpcClient.Close()

	// Adapter is just a wrapper to create Service client from protogen file passing it the grpc client created earlier
	helloAdapter, err := adapter.NewHelloAdapter(grpcClient)
	if err != nil {
		log.Fatalf("Error while creating HelloAdapter :%v", err)
	}

	// runSayHello(*helloAdapter, "Bruce Banner")

	// helloAdapter.SayHelloServerStream(context.Background(), "Scarlet Witch")

	// helloAdapter.SayHelloClientStream(context.Background(), []string{"Vallabh", "Ashish", "Steve", "Somnath"})

	// Call the gRPC server method from adapter wrapper by making use of the service client created earlier
	helloAdapter.SayHelloContinuous(context.Background(), []string{"Superman", "Joker", "Batman", "Aquaman", "Flash"})

	bAdapter, err := bankadapter.NewBankAdapter(grpcClient)
	if err != nil {
		log.Fatalf("Error while creating BankAdapter :%v", err)
	}

	runGetCurrentBalance(bAdapter, "7835697001xx")
	// runFetchExchangeRates(bAdapter, "USD", "INR")
	//runSummarizeTransactions(bAdapter, "7835697002", 5)
	// runTransferMultiple(bAdapter, "7835697001", "7835697002", 10)
}

func runSayHello(adapter adapter.HelloAdapter, name string) {
	greet, err := adapter.SayHello(context.Background(), name)
	if err != nil {
		log.Fatalf("Error while creating HelloAdapter :%v", err)
	}
	log.Println(greet.Greet)
}

func runGetCurrentBalance(adapter bankadapter.BankAdapter, acct string) {
	bal, err := adapter.GetCurrentBalance(context.Background(), acct)
	if err != nil {
		log.Fatalln("Failed to call GetCurrentBalance: ", err)
	}
	log.Println("Current balance: ", bal)
}

func runFetchExchangeRates(adapter bankadapter.BankAdapter, fromAcct, toAcct string) {
	adapter.FetchExchangeRates(context.Background(), fromAcct, toAcct)
}

func runSummarizeTransactions(adapter bankadapter.BankAdapter, acct string, numDummyTransactions int) {

	var txs []bank.Transaction
	for i := 1; i <= numDummyTransactions; i++ {
		ttype := bank.TransactionTypeIn

		if i%3 == 0 {
			ttype = bank.TransactionTypeOut
		}

		t := bank.Transaction{
			Amount:          float64(rand.Intn(500) + 10),
			TransactionType: ttype,
			Notes:           fmt.Sprint("Dummy transaction:%v", i),
		}

		txs = append(txs, t)
	}

	adapter.SummarizeTransactions(context.Background(), acct, txs)
}

func runTransferMultiple(adapter bankadapter.BankAdapter, fromAcct, toAcct string, numDummyTransactions int) {

	var trf []bank.TransferTransaction

	for i := 1; i <= numDummyTransactions; i++ {
		tr := bank.TransferTransaction{
			FromAccountNumber: fromAcct,
			ToAccountNumber:   toAcct,
			Currency:          "USD",
			Amount:            float64(rand.Intn(200) + 5),
		}

		trf = append(trf, tr)
	}

	adapter.TransferMultiple(context.Background(), trf)

}
