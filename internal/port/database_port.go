package port

import (
	"github.com/google/uuid"
	db "github.com/seminhnva/my-grpc-go-server/internal/adapter/database"
)

type DummyDatabasePort interface {
	Save(data *db.DummyOrm) (uuid.UUID, error)
	GetByUuid(uuid *uuid.UUID) (db.DummyOrm, error)
}

type BankDatabasePort interface {
	GetBankAccountByAccountNumber(accountNumber string) (db.BankAccountOrm, error)
	InsertDummyExchangeRate(r db.BankExchangeRateOrm) error
	FetchExchangeRate() (db.BankExchangeRateOrm, error)
	CreateTransaction(acc db.BankAccountOrm, transaction db.BankTransactionOrm) (uuid.UUID, error)
	CreateTranfer(transfer db.BankTransferOrm) (uuid.UUID, error)
	CreateTransferTransactionPair(fromAcc db.BankAccountOrm, toAcc db.BankAccountOrm, fromTransaction db.BankTransactionOrm, toTransaction db.BankTransactionOrm) (bool, error)
	UpdateTransferStatus(transfer db.BankTransferOrm, status bool) error
}
