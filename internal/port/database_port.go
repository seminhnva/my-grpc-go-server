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
	GetCurrentBalance(accountNumber string) (db.BankAccountOrm, error)
}
