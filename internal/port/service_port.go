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
	GetLatestExchaneRate(fromCurrency, toCurrency string) (float64, error)
	CreateTransaction(accountNumber string, t dbank.Transaction) (uuid.UUID, error)
	CalculateTransactionSumary(tcur *dbank.TransactionSummary, t dbank.Transaction) error
	Transfer(tt dbank.TransferTransaction) (uuid.UUID, bool, error)
}
