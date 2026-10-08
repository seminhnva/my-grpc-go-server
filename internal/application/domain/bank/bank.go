package bank

import "time"

type ExchangeRate struct {
	FromCurrency       string
	ToCurrency         string
	Rate               float64
	ValidFromTimeStamp time.Time
	ValidToTimeStamp   time.Time
}
