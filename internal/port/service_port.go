package port

import (
	"github.com/google/uuid"
	dbank "github.com/seminhnva/my-grpc-go-server/internal/application/domain/bank"
)

type HelloSerivcePort interface {
	GenerateGreet(name string) string
}

type BankServicePort interface {
	FindCurrentBalance(account string) float64
	CreateDummyExchangeRate(r dbank.ExchangeRate)
	GetLatestExchaneRate() float64
	CreateTransaction(accountNumber string, t dbank.Transaction) (uuid.UUID, error)
	CalculateTransactionSumary(tcur *dbank.TransactionSummary, t dbank.Transaction) error
	TransferMultiple(transfer dbank.Transfer) (uuid.UUID, error)
}
