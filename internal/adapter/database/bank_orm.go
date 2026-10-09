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
	TransactionUUID      uuid.UUID `gorm:"primaryKey"`
	AccountUUID          uuid.UUID `gorm:"type:uuid;not null"`
	TransactionTimestamp time.Time
	Amount               float64
	TransactionType      string
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

type BankTransferOrm struct {
	TransferUUID      uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"transfer_uuid"`
	FromAccountUUID   uuid.UUID  `gorm:"type:uuid" json:"from_account_uuid"`
	ToAccountUUID     uuid.UUID  `gorm:"type:uuid" json:"to_account_uuid"`
	Currency          string     `gorm:"type:varchar(5);not null" json:"currency"`
	Amount            float64    `gorm:"type:numeric(15,2);not null" json:"amount"`
	TransferTimestamp time.Time  `gorm:"type:timestamptz;not null" json:"transfer_timestamp"`
	TransferSuccess   bool       `gorm:"not null;default:false" json:"transfer_success"`
	CreatedAt         *time.Time `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt         *time.Time `gorm:"type:timestamptz" json:"updated_at"`
}

func (BankTransferOrm) TableName() string {
	return "bank_transfers"
}
