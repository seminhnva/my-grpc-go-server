package grpc

import (
	"context"
	"time"

	"github.com/seminhnva/my-grpc-proto/protogen/go/bank"
	"google.golang.org/grpc"
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
