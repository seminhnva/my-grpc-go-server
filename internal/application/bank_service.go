package application

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	db "github.com/seminhnva/my-grpc-go-server/internal/adapter/database"
	"github.com/seminhnva/my-grpc-go-server/internal/application/domain/bank"
	dbank "github.com/seminhnva/my-grpc-go-server/internal/application/domain/bank"
	"github.com/seminhnva/my-grpc-go-server/internal/port"
)

type BankService struct {
	port port.BankDatabasePort
}

func NewBankService(db port.BankDatabasePort) *BankService {
	return &BankService{
		port: db,
	}
}

func (bs *BankService) FindCurrentBalance(accountName string) float64 {
	accountInfo, err := bs.port.GetCurrentBalance(accountName)
	if err != nil {
		log.Println("Cant find accountName", err)
	}
	return accountInfo.CurrentBalance
}
func (bs *BankService) CreateDummyExchangeRate(
	r dbank.ExchangeRate,
) {
	now := time.Now()

	exchangeRateOrm := db.BankExchangeRateOrm{
		ExchangeRateUUID:   uuid.New(),
		FromCurrency:       r.FromCurrency,
		ToCurrency:         r.ToCurrency,
		Rate:               r.Rate,
		ValidFromTimestamp: now,
		ValidToTimestamp:   now.Add(1 * time.Hour),
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := bs.port.InsertDummyExchangeRate(exchangeRateOrm); err != nil {
		log.Println("Insert exchange rate failed:", err)
	}
}

func (bs *BankService) GetLatestExchaneRate() float64 {
	rateInfo, err := bs.port.FetchExchangeRate()
	if err != nil {
		log.Println("Cant get exchange rate", err)
		return 0
	}
	return rateInfo.Rate
}

func (bs *BankService) CreateTransaction(accountNumber string, t dbank.Transaction) (uuid.UUID, error) {
	accountInfo, err := bs.port.GetCurrentBalance(accountNumber)
	if err != nil {
		log.Printf("Can't create transaction for %v : %v\n", accountNumber, err)
		return uuid.Nil, fmt.Errorf("can't find account number %v : %v", accountNumber, err.Error())
	}

	if t.TransactionType == bank.TransactionTypeOut && accountInfo.CurrentBalance < t.Amount {
		return accountInfo.AccountUUID, fmt.Errorf(
			"insufficient account balance %v for [out] transaction amount %v",
			accountInfo.CurrentBalance, t.Amount,
		)
	}
	newUuid := uuid.New()
	now := time.Now()
	transOrm := db.BankTransactionOrm{
		TransactionUUID:      newUuid,
		AccountUUID:          accountInfo.AccountUUID,
		TransactionTimestamp: time.Now(),
		Amount:               t.Amount,
		Notes:                t.Notes,
		TransactionType:      t.TransactionType,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	savedUuid, err := bs.port.CreateTransaction(accountInfo, transOrm)

	return savedUuid, err
}

func (bs *BankService) CalculateTransactionSumary(tcur *dbank.TransactionSummary, t dbank.Transaction) error {
	switch t.TransactionType {
	case dbank.TransactionTypeIn:
		tcur.SumIn += t.Amount
	case dbank.TransactionTypeOut:
		tcur.SumOut += t.Amount
	default:
		return fmt.Errorf("Unknown transaction typ: %s", t.TransactionType)
	}
	tcur.SumTotal = tcur.SumIn - tcur.SumOut
	return nil
}

func (bs *BankService) TransferMultiple(transfer dbank.Transfer) (uuid.UUID, error) {
	fromAccInfo, err := bs.port.GetCurrentBalance(transfer.FromAccountNumber)
	if err != nil {
		log.Printf("Can't create transfer for %v : %v\n", fromAccInfo.AccountName, err)
		return uuid.Nil, fmt.Errorf("can't find account number %v : %v", fromAccInfo.AccountName, err.Error())
	}
	toAccInfo, err := bs.port.GetCurrentBalance(transfer.ToAccountNumber)
	if err != nil {
		log.Printf("Can't create transfer for %v : %v\n", toAccInfo.AccountName, err)
		return uuid.Nil, fmt.Errorf("can't find account number %v : %v", toAccInfo.AccountName, err.Error())
	}

	tranfer := db.BankTransferOrm{
		TransferUUID:    uuid.New(),
		FromAccountUUID: fromAccInfo.AccountUUID,
		ToAccountUUID:   toAccInfo.AccountUUID,
		Currency:        transfer.Currency,
		Amount:          transfer.Amount,
	}

	tranferID, err := bs.port.TransferMultiple(fromAccInfo, toAccInfo, tranfer)
	if err != nil {
		log.Printf("cant tranfer multiple: %v", err)
	}
	return tranferID, nil

}
