package database

import (
	"time"

	"github.com/google/uuid"
)

type BankAccountOrm struct {
	AccountUUID    uuid.UUID `gorm:"primaryKey"`
	AccountNumber  string
	AccountName    string
	Currency       string
	CurrentBalance float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (BankAccountOrm) TableName() string {
	return "bank_accounts"
}

type BankTransactionOrm struct {
	TransacctionUUID     uuid.UUID `gorm:"primaryKey"`
	AccountUUID          uuid.UUID `gorm:"type:uuid;not null"`
	TransactionTimeStamp time.Time
	Amount               float64
	TransactionType      float64
	Notes                string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Account              BankAccountOrm `gorm:"foreignKey:AccountUUID;references:AccountUUID"`
}

func (BankTransactionOrm) TableName() string {
	return "bank_transactions"
}
