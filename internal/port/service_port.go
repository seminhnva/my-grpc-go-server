package port

import (
	"github.com/google/uuid"
	dbank "github.com/seminhnva/my-grpc-go-server/internal/application/domain/bank"
)

type HelloSerivcePort interface {
	GenerateGreet(name string) string
}

type BankServicePort interface {
	FindCurrentBalance(account string) (float64, error)
	CreateDummyExchangeRate(r dbank.ExchangeRate) error
	GetLatestExchaneRate(fromCurrency, toCurrency string) (dbank.ExchangeRate, error)
	CreateTransaction(accountNumber string, t dbank.Transaction) (uuid.UUID, error)
	CalculateTransactionSumary(tcur *dbank.TransactionSummary, t dbank.Transaction) error
	Transfer(tt dbank.TransferTransaction) (uuid.UUID, bool, error)
}

type ResiliencyServicePort interface {
	GenerateResiliency(minDelaySecond int32, maxDelaySecond int32, statusCode []uint32) (string, uint32)
}
