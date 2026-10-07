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

type BankExchangeRateOrm struct {
	ExchangeRateUUID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	FromCurrency       string    `gorm:"type:varchar(5);not null"`
	ToCurrency         string    `gorm:"type:varchar(5);not null"`
	Rate               float64   `gorm:"type:numeric(20,10);not null"`
	ValidFromTimestamp time.Time `gorm:"not null"`
	ValidToTimestamp   time.Time `gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (BankExchangeRateOrm) TableName() string {
	return "bank_exchange_rates"
}
