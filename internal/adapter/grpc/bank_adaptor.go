package grpc

import (
	"context"
	"time"

	"github.com/seminhnva/my-grpc-proto/protogen/go/bank"
	"google.golang.org/genproto/googleapis/type/date"
)

func (a *GrpcAdapter) GetCurrentBalance(ctx context.Context, req *bank.CurrentBalanceRequest) (*bank.CurrrentBalanceResponse, error) {
	now := time.Now()
	amount := a.bankService.FindCurrentBalance(req.AccountNumber)

	return &bank.CurrrentBalanceResponse{
		CurrentDate: &date.Date{
			Year:  int32(now.Year()),
			Month: int32(now.Month()),
			Day:   int32(now.Day()),
		},
		Amount: amount,
	}, nil
}
