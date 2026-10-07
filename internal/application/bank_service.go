package application

import (
	"log"

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
