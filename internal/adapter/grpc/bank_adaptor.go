package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	dbank "github.com/seminhnva/my-grpc-go-server/internal/application/domain/bank"
	"github.com/seminhnva/my-grpc-proto/protogen/go/bank"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *GrpcAdapter) GetCurrentBalance(ctx context.Context, req *bank.CurrentBalanceRequest) (*bank.CurrrentBalanceResponse, error) {
	currentBalance := a.bankService.FindCurrentBalance(req.AccountNumber)
	return &bank.CurrrentBalanceResponse{
		CurrentBalance: currentBalance,
	}, nil
}

func (a *GrpcAdapter) FetchExchangeRates(req *bank.ExchangeRateRequest, stream grpc.ServerStreamingServer[bank.ExchangeRateResponse]) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			rate := a.bankService.GetLatestExchaneRate()
			if err := stream.Send(&bank.ExchangeRateResponse{
				Rate: rate,
			}); err != nil {
				return err
			}

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}
func toTime(dt *datetime.DateTime) (time.Time, error) {
	if dt == nil {
		now := time.Now()

		dt = &datetime.DateTime{
			Year:    int32(now.Year()),
			Month:   int32(now.Month()),
			Day:     int32(now.Day()),
			Hours:   int32(now.Hour()),
			Minutes: int32(now.Minute()),
			Seconds: int32(now.Second()),
			Nanos:   int32(now.Nanosecond()),
		}
	}

	res := time.Date(
		int(dt.Year),
		time.Month(dt.Month),
		int(dt.Day),
		int(dt.Hours),
		int(dt.Minutes),
		int(dt.Seconds),
		int(dt.Nanos),
		time.UTC,
	)

	return res, nil
}

func (a *GrpcAdapter) SummarizeTransactions(stream grpc.ClientStreamingServer[bank.Transaction, bank.TransactionSummary]) error {
	tSumary := dbank.TransactionSummary{
		SummaryOnDate: time.Now(),
		SumIn:         0,
		SumOut:        0,
		SumTotal:      0,
	}
	accountNumber := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&bank.TransactionSummary{
				AccountNumber: accountNumber,
				SumAmountIn:   tSumary.SumIn,
				SumAmountOut:  tSumary.SumOut,
				SumTotal:      tSumary.SumTotal,
				TransactionDate: &date.Date{
					Year:  int32(tSumary.SummaryOnDate.Year()),
					Month: int32(tSumary.SummaryOnDate.Month()),
					Day:   int32(tSumary.SummaryOnDate.Day()),
				},
			})
		}
		if err != nil {
			log.Fatalln("Error while reading from client :", err)

		}
		accountNumber = req.AccountNumber
		ts, err := toTime(req.Timestamp)
		if err != nil {
			log.Fatalf("Error while parsing timestamp %v : %v", req.Timestamp, err)
		}
		ttype := ""
		switch req.Type {
		case bank.TransactionType_TRANSACTION_TYPE_IN:
			ttype = dbank.TransactionTypeIn
		case bank.TransactionType_TRANSACTION_TYPE_OUT:
			ttype = dbank.TransactionTypeOut
		default:
			ttype = dbank.TransactionTypeUnknown
		}

		transaction := dbank.Transaction{
			Amount:          req.Amount,
			Timestamp:       ts,
			TransactionType: ttype,
			Notes:           req.Notes,
		}
		accountUuid, err := a.bankService.CreateTransaction(accountNumber, transaction)

		if err != nil && accountUuid == uuid.Nil {
			s := status.New(codes.InvalidArgument, err.Error())
			s, _ = s.WithDetails(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "account_number",
						Description: "Invalid account number",
					},
				},
			})

			return s.Err()
		} else if err != nil && accountUuid != uuid.Nil {
			s := status.New(codes.InvalidArgument, err.Error())
			s, _ = s.WithDetails(&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "amount",
						Description: fmt.Sprintf("Requested amount %v exceed available balance", req.Amount),
					},
				},
			})

			return s.Err()
		}
		if err != nil {
			log.Println("Error while creating transaction :", err)
		}

		err = a.bankService.CalculateTransactionSumary(&tSumary, transaction)

		if err != nil {
			log.Println("Error while calculate sumarize :", err)
		}
	}
}

func (a *GrpcAdapter) TransferMultiple(stream grpc.BidiStreamingServer[bank.TransferRequest, bank.TransferResponse]) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		tranfer := dbank.Transfer{
			FromAccountNumber: req.FromAccountNumber,
			ToAccountNumber:   req.ToAccountNumber,
			Currency:          req.Currency,
			Amount:            req.Amount,
		}
		_, err = a.bankService.TransferMultiple(tranfer)
		if err != nil {
			log.Fatalln("Error while say HelloContinious", err)
		}
		if err := stream.Send(&bank.TransferResponse{
			FromAccountNumber: tranfer.FromAccountNumber,
			ToAccountNumber:   tranfer.ToAccountNumber,
			Currency:          tranfer.Currency,
			Amount:            tranfer.Amount,
			Status:            bank.TransferStatus_TRANSFER_STATUS_SUCCESS,
		}); err != nil {
			log.Fatalln("Error while say HelloContinious", err)
		}
	}
}
