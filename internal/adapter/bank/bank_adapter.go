package bank

import (
	"context"
	"io"
	"log"

	protogenbank "github.com/VallabhSLEPAM/go-with-grpc/protogen/go/bank"
	"github.com/VallabhSLEPAM/grpc-client/internal/application/domain/bank"
	"github.com/VallabhSLEPAM/grpc-client/internal/port"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BankAdapter struct {
	bankClient port.BankClientPort
}

func NewBankAdapter(conn *grpc.ClientConn) (BankAdapter, error) {
	client := protogenbank.NewBankServiceClient(conn)
	return BankAdapter{
		bankClient: client,
	}, nil
}

func (adapter BankAdapter) GetCurrentBalance(ctx context.Context, acctNumber string) (*protogenbank.CurrentBalanceResponse, error) {

	bankRequest := protogenbank.CurrentBalanceRequest{
		AccountNumber: acctNumber,
	}

	bal, err := adapter.bankClient.GetCurrentBalance(ctx, &bankRequest)
	if err != nil {
		st, _ := status.FromError(err)
		log.Fatalf("[FATAL] Error on GetCurrentBalance: ", st)
	}
	return bal, nil
}

func (adapter BankAdapter) FetchExchangeRates(ctx context.Context, fromCurr, toCurr string) {

	exchangeRateRequest := protogenbank.ExchangeRateRequest{
		FromCurrency: fromCurr,
		ToCurrency:   toCurr,
	}

	exchangeRateStream, err := adapter.bankClient.FetchExchangeRates(ctx, &exchangeRateRequest)
	if err != nil {
		st, _ := status.FromError(err)
		log.Fatalf("[FATAL] Error fetching the exchange rates from %v to %v with error %v\n", fromCurr, toCurr, st)
	}

	for {
		rates, err := exchangeRateStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() == codes.InvalidArgument {
				log.Fatalln("[FATAL] Error on FetchExchangeRates: \n", st.Message())

			}
		}

		log.Printf("Rates at %v from %v to %v: %v\n ", rates.Timestamp, fromCurr, toCurr, rates.Rate)
	}
}

func (adapter BankAdapter) SummarizeTransactions(ctx context.Context, acct string, txs []bank.Transaction) {

	txStream, err := adapter.bankClient.SummarizeTransactions(ctx)
	if err != nil {
		log.Fatalln("[FATAL] Error on SummarizeTransactions: ", err)
	}

	for _, tx := range txs {
		ttype := protogenbank.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
		if tx.TransactionType == bank.TransactionTypeIn {
			ttype = protogenbank.TransactionType_TRANSACTION_TYPE_IN
		} else if tx.TransactionType == bank.TransactionTypeOut {
			ttype = protogenbank.TransactionType_TRANSACTION_TYPE_OUT
		}

		transactionRequest := protogenbank.Transaction{
			AccountNumber: acct,
			Type:          ttype,
			Amount:        tx.Amount,
			Notes:         tx.Notes,
		}
		txStream.Send(&transactionRequest)
	}

	summary, err := txStream.CloseAndRecv()
	if err != nil {
		st, _ := status.FromError(err)
		log.Fatalf("[FATAL] Error receiving server stream response: %v\n", st)
	}
	log.Println("Summary details: ", summary)
}

func (adapter BankAdapter) TransferMultiple(ctx context.Context, trf []bank.TransferTransaction) {

	clienttxsStream, err := adapter.bankClient.TransferMultiple(ctx)
	if err != nil {
		log.Fatalf("[FATAL] Error in TransferMultiple: %v\n", err)
	}
	trfChan := make(chan chan struct{})

	go func() {
		for _, tr := range trf {
			req := &protogenbank.TransferRequest{
				FromAccountNumber: tr.FromAccountNumber,
				ToAccountNumber:   tr.ToAccountNumber,
				Current:           tr.Currency,
				Amount:            float32(tr.Amount),
			}
			clienttxsStream.Send(req)
		}
		clienttxsStream.CloseSend()
	}()

	go func() {
		for {
			res, err := clienttxsStream.CloseAndRecv()
			if err == io.EOF {
				break
			}
			if err != nil {
				handleTransferErrorGrpc(err)
				break
			} else {
				log.Printf("Transfer status %v on %v", res.Status, res.Timestamp)
			}
		}
		close(trfChan)
	}()

	<-trfChan

}

func handleTransferErrorGrpc(err error) {
	st := status.Convert(err)

	log.Printf("Error %v on TransferMultiple : %v", st.Code(), st.Message())

	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.PreconditionFailure:
			for _, violation := range t.GetViolations() {
				log.Println("[VIOLATION]", violation)
			}
		case *errdetails.ErrorInfo:
			log.Printf("Error on : %v, with reason %v\n", t.Domain, t.Reason)
			for k, v := range t.GetMetadata() {
				log.Printf("  %v : %v\n", k, v)
			}
		}
	}
}
