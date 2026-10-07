package grpc

import (
	"context"

	"github.com/seminhnva/my-grpc-proto/protogen/go/bank"
)

func (a *GrpcAdapter) GetCurrentBalance(ctx context.Context, req *bank.CurrentBalanceRequest) (*bank.CurrrentBalanceResponse, error) {
	currentBalance := a.bankService.FindCurrentBalance(req.AccountNumber)
	return &bank.CurrrentBalanceResponse{
		CurrentBalance: currentBalance,
	}, nil
}
