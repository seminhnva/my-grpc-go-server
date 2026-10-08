package database

import (
	"log"
	"time"
)

func (a *DatabaseAdapter) GetCurrentBalance(accountNumber string) (BankAccountOrm, error) {
	var res BankAccountOrm
	if err := a.db.First(&res, "account_number = ?", accountNumber).Error; err != nil {
		log.Println("Cant get data ", err)
		return res, err
	}
	return res, nil
}

func (a *DatabaseAdapter) InsertDummyExchangeRate(r BankExchangeRateOrm) error {

	if err := a.db.Create(r).Error; err != nil {
		log.Println("Insert exchange rate failed:", err)
		return err
	}
	log.Printf(
		"Inserted exchange rate %s -> %s = %.10f",
		r.FromCurrency,
		r.ToCurrency,
		r.Rate,
	)
	return nil
}
func (a *DatabaseAdapter) FetchExchangeRate() (BankExchangeRateOrm, error) {
	var rate BankExchangeRateOrm
	err := a.db.
		Where(
			"from_currency = ? AND to_currency = ? AND valid_from_timestamp <= ? AND valid_to_timestamp > ?",
			"USD",
			"VND",
			time.Now(),
			time.Now(),
		).
		Order("valid_from_timestamp DESC").
		First(&rate).Error
	if err != nil {
		log.Println("Fail to fetch rate")
		return rate, err
	}
	return rate, err
}
