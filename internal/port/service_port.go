package port

import (
	dbank "github.com/seminhnva/my-grpc-go-server/internal/application/domain/bank"
)

type HelloSerivcePort interface {
	GenerateGreet(name string) string
}

type BankServicePort interface {
	FindCurrentBalance(account string) float64
	CreateDummyExchangeRate(r dbank.ExchangeRate)
	GetLatestExchaneRate() float64
}
