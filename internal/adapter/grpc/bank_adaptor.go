package grpc

import (
	"context"
	"errors"
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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *GrpcAdapter) GetCurrentBalance(ctx context.Context, req *bank.CurrentBalanceRequest) (*bank.CurrrentBalanceResponse, error) {
	currentBalance, err := a.bankService.FindCurrentBalance(req.AccountNumber)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "account %v not found", req.AccountNumber)
	}
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
			rate, err := a.bankService.GetLatestExchaneRate(req.FromCurrency, req.ToCurrency)

			if err != nil {
				if errors.Is(err, dbank.ErrExchangeRateNotFound) {
					return status.Errorf(
						codes.NotFound,
						"exchange rate record not found for %s -> %s",
						req.FromCurrency,
						req.ToCurrency,
					)
				}
				return status.Errorf(codes.Internal, "failed to fetch exchange rate: %v", err)
			}
			if err := stream.Send(&bank.ExchangeRateResponse{
				FromCurrency:         rate.FromCurrency,
				ToCurrency:           rate.ToCurrency,
				Rate:                 rate.Rate,
				TransactionTimestamp: timestamppb.New(rate.ValidFromTimeStamp),
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

func currentDatetime() *datetime.DateTime {
	now := time.Now()

	return &datetime.DateTime{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		Day:        int32(now.Day()),
		Hours:      int32(now.Hour()),
		Minutes:    int32(now.Minute()),
		Seconds:    int32(now.Second()),
		Nanos:      int32(now.Nanosecond()),
		TimeOffset: &datetime.DateTime_UtcOffset{},
	}
}

func (a *GrpcAdapter) TransferMultiple(stream grpc.BidiStreamingServer[bank.TransferRequest, bank.TransferResponse]) error {
	context := stream.Context()
	for {
		select {
		case <-context.Done():
			log.Println("Client cancel stream")
			return context.Err()
		default:
			req, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				// log.Fatalf("Error while reading from client %v:", err)
				return err
			}
			tt := dbank.TransferTransaction{
				FromAccountNumber: req.FromAccountNumber,
				ToAccountNumber:   req.ToAccountNumber,
				Currency:          req.Currency,
				Amount:            req.Amount,
			}

			_, transferSuccess, err := a.bankService.Transfer(tt)
			if err != nil {
				return buildTransferErrorStatusGrpc(err, req)
			}

			res := bank.TransferResponse{
				FromAccountNumber: req.FromAccountNumber,
				ToAccountNumber:   req.ToAccountNumber,
				Currency:          req.Currency,
				Amount:            req.Amount,
				Timestamp:         currentDatetime(),
			}
			if transferSuccess {
				res.Status = bank.TransferStatus_TRANSFER_STATUS_SUCCESS
			} else {
				res.Status = bank.TransferStatus_TRANSFER_STATUS_FAILED
			}
			err = stream.Send(&res)

			if err != nil {
				log.Fatalln("Error while sending response to client :", err)
			}
		}

	}
}

func buildTransferErrorStatusGrpc(err error, req *bank.TransferRequest) error {
	switch {
	case errors.Is(err, dbank.ErrTransferSourceAccountNotFound):
		s := status.New(codes.FailedPrecondition, err.Error())
		s, _ = s.WithDetails(&errdetails.PreconditionFailure{
			Violations: []*errdetails.PreconditionFailure_Violation{
				{
					Type:        "INVALID_ACCOUNT",
					Subject:     "Source account not found",
					Description: fmt.Sprintf("source account (from %v) not found", req.FromAccountNumber),
				},
			},
		})

		return s.Err()
	case errors.Is(err, dbank.ErrTransferDestinationAccountNotFound):
		s := status.New(codes.FailedPrecondition, err.Error())
		s, _ = s.WithDetails(&errdetails.PreconditionFailure{
			Violations: []*errdetails.PreconditionFailure_Violation{
				{
					Type:        "INVALID_ACCOUNT",
					Subject:     "Destination account not found",
					Description: fmt.Sprintf("destination account (to %v) not found", req.ToAccountNumber),
				},
			},
		})

		return s.Err()
	case errors.Is(err, dbank.ErrTransferRecordFailed):
		s := status.New(codes.Internal, err.Error())
		s, _ = s.WithDetails(&errdetails.Help{
			Links: []*errdetails.Help_Link{
				{
					Url:         "my-bank-website.com/faq",
					Description: "Bank FAQ",
				},
			},
		})

		return s.Err()
	case errors.Is(err, dbank.ErrTransferTransactionFair):
		s := status.New(codes.InvalidArgument, err.Error())
		s, _ = s.WithDetails(&errdetails.ErrorInfo{
			Domain: "my-bank-website.com",
			Reason: "TRANSACTION_PAIR_FAILED",
			Metadata: map[string]string{
				"from_account": req.FromAccountNumber,
				"to_account":   req.ToAccountNumber,
				"currency":     req.Currency,
				"amount":       fmt.Sprintf("%f", req.Amount),
			},
		})

		return s.Err()
	default:
		s := status.New(codes.Unknown, err.Error())
		return s.Err()
	}
}
