package database

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/seminhnva/my-grpc-go-server/internal/application/domain/bank"
	"gorm.io/gorm"
)

func (a *DatabaseAdapter) GetBankAccountByAccountNumber(accountNumber string) (BankAccountOrm, error) {
	var res BankAccountOrm
	if err := a.db.First(&res, "account_number = ?", accountNumber).Error; err != nil {
		log.Println("Cant get data ", err)
		return res, err
	}
	return res, nil
}

func (a *DatabaseAdapter) InsertDummyExchangeRate(r BankExchangeRateOrm) error {
	if err := a.db.Create(&r).Error; err != nil {
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
func (a *DatabaseAdapter) FetchExchangeRate(fromCurrency, toCurrency string) (BankExchangeRateOrm, error) {
	var rate BankExchangeRateOrm
	err := a.db.
		Where(
			"from_currency = ? AND to_currency = ? AND valid_from_timestamp <= ? AND valid_to_timestamp > ?",
			fromCurrency,
			toCurrency,
			time.Now(),
			time.Now(),
		).
		Order("valid_from_timestamp DESC").
		First(&rate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rate, bank.ErrExchangeRateNotFound
		}
		return rate, err
	}
	return rate, nil
}
func (a *DatabaseAdapter) CreateTransaction(acc BankAccountOrm, transaction BankTransactionOrm) (uuid.UUID, error) {
	tx := a.db.Begin()
	transAmount := transaction.Amount
	if transaction.TransactionType == bank.TransactionTypeOut {
		transAmount = -1 * transaction.Amount
	}
	newAmount := acc.CurrentBalance + transAmount
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err := tx.Model(&acc).Updates(
		map[string]interface{}{
			"current_balance": newAmount,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return uuid.Nil, err
	}
	return transaction.TransactionUUID, nil
}

func (a *DatabaseAdapter) CreateTranfer(transfer BankTransferOrm) (uuid.UUID, error) {
	if err := a.db.Create(&transfer).Error; err != nil {
		return uuid.Nil, err
	}
	return uuid.Nil, nil
}

func (a *DatabaseAdapter) CreateTransferTransactionPair(fromAcc BankAccountOrm, toAcc BankAccountOrm, fromTransaction BankTransactionOrm, toTransaction BankTransactionOrm) (bool, error) {
	tx := a.db.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}
	defer tx.Rollback()

	if fromAcc.AccountUUID == toAcc.AccountUUID {
		return false, errors.New("cannot transfer to the same account")
	}

	if err := tx.Create(&fromTransaction).Error; err != nil {
		return false, err
	}
	if err := tx.Create(&toTransaction).Error; err != nil {
		return false, err
	}

	result := tx.Model(&BankAccountOrm{}).
		Where(
			"account_uuid = ? AND current_balance >= ?",
			fromAcc.AccountUUID,
			fromTransaction.Amount,
		).
		Updates(map[string]interface{}{
			"current_balance": gorm.Expr("current_balance - ?", fromTransaction.Amount),
			"updated_at":      time.Now(),
		})

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected != 1 {
		return false, bank.ErrTransferTransactionFair
	}

	result = tx.Model(&BankAccountOrm{}).
		Where("account_uuid = ?", toAcc.AccountUUID).
		Updates(map[string]interface{}{
			"current_balance": gorm.Expr("current_balance + ?", toTransaction.Amount),
			"updated_at":      time.Now(),
		})

	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected != 1 {
		return false, bank.ErrTransferDestinationAccountNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return false, err
	}
	return true, nil
}
func (a *DatabaseAdapter) UpdateTransferStatus(transfer BankTransferOrm, status bool) error {
	if err := a.db.Model(&transfer).Updates(
		map[string]interface{}{
			"transfer_success": status,
			"updated_at":       time.Now(),
		},
	).Error; err != nil {
		return err
	}
	return nil
}
