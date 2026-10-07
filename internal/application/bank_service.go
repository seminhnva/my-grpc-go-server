package application

import (
	"log"
	"time"

	"github.com/seminhnva/my-grpc-go-server/internal/port"
)

type BankSerivce struct {
	port port.BankDatabasePort
}

func NewBankService(db port.BankDatabasePort) *BankSerivce {
	return &BankSerivce{
		port: db,
	}
}

func (bs *BankSerivce) FindCurrentBalance(accountName string) float64 {
	accountInfo, err := bs.port.GetCurrentBalance(accountName)
	if err != nil {
		log.Println("Cant find accountName", err)
	}
	return accountInfo.CurrentBalance
}
func (bs *BankSerivce) StartDummyExchangeRateInjector(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	bs.port.InsertDummyExchangeRate()

	for range ticker.C {
		bs.port.InsertDummyExchangeRate()
	}
}

func (bs *BankSerivce) GetLatestExchaneRate() float64 {
	rateInfo, err := bs.port.FetchExchangeRate()
	if err != nil {
		log.Println("Cant get exchange rate", err)
	}
	return rateInfo.Rate
}
