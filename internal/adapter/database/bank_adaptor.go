package database

import (
	"log"
)

func (a *DatabaseAdapter) GetCurrentBalance(accountNumber string) (BankAccountOrm, error) {
	var res BankAccountOrm
	if err := a.db.First(&res, "account_number = ?", accountNumber).Error; err != nil {
		log.Println("Cant get data ", err)
		return res, err
	}
	return res, nil
}
