package application

import (
	"log"
	"time"

	"github.com/google/uuid"
	db "github.com/seminhnva/my-grpc-go-server/internal/adapter/database"
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
