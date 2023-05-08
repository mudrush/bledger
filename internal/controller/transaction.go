package controller

import (
	"encoding/json"
	"errors"

	"github.com/partyscript/bledger/internal/cache"
	"github.com/partyscript/bledger/internal/common"
	"github.com/partyscript/bledger/internal/config"
	"github.com/partyscript/bledger/internal/db"
	"github.com/partyscript/bledger/internal/model"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TransactionsController is the struct that the constructor implements
type TransactionsController struct {
	logger *zap.SugaredLogger
	cfg    *config.GlobalConfig
	cache  *cache.Manager
	db     *db.Manager
}

// NewTransactionsController initializes a TransactionsController instance
func NewTransactionsController(
	logger *zap.SugaredLogger,
	cfg *config.GlobalConfig,
	cache *cache.Manager,
	db *db.Manager,
) TransactionsController {
	return TransactionsController{
		logger: logger,
		cfg:    cfg,
		cache:  cache,
		db:     db,
	}
}

// CreatePendingTransaction creates a new transaction with a pending state
func (tc *TransactionsController) CreatePendingTransaction(txReq model.CreateTransactionRequest) (*model.Transaction, error) {
	var acct model.Account

	amtToStore, err := json.Marshal(txReq.Money)
	if err != nil {
		return nil, err
	}

	args := model.Transaction{
		AccountID: txReq.AccountID,
		Money:     datatypes.JSON(amtToStore),
		Direction: txReq.Direction,
		Memo:      txReq.Memo,
		State:     model.TransactionStatePending,
		Version:   0,
	}

	find := tc.db.Gorm.First(&acct, &model.Account{ID: args.AccountID})
	if find.Error != nil {
		return nil, errors.New("account not found")
	}

	if !isValidCurrency(args.Money, acct.Balance) {
		args.State = model.TransactionStateFailed
		args.ErrorReason = "invalid currency"
	} else {
		args.State = model.TransactionStatePending
	}

	// Calculate if we should automatically fail this
	err = calcAccountBalanceChange(&args, &acct)
	if err != nil {
		args.State = model.TransactionStateFailed
		args.ErrorReason = err.Error()
	}

	create := tc.db.Gorm.Create(&args)
	if create.Error != nil {
		return nil, create.Error
	}

	return &args, nil
}

// ExecutePendingTransaction executes a pending transaction
func (tc *TransactionsController) ExecutePendingTransaction(txID string) (*model.Transaction, error) {
	var foundTx model.Transaction
	var acct model.Account

	// start tx pipeline
	tx := tc.db.Gorm.Begin()

	// lock the tx row
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&foundTx, &model.Transaction{ID: txID}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// lock the account row
	err = tx.Set("gorm:query_option", "FOR UPDATE").First(&acct, &model.Account{ID: foundTx.AccountID}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = executeAccountChanges(tx, &foundTx, &acct, false)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &foundTx, nil
}

// GetTransaction returns a transaction by id
func (tc *TransactionsController) GetTransaction(txID string) (*model.Transaction, error) {
	tx := new(model.Transaction)

	find := tc.db.Gorm.First(&tx, "id = ?", txID)

	if find.Error != nil {
		return nil, find.Error
	}

	return tx, nil
}

// CreateTransaction creates a new transaction with immediate clearance
func (tc *TransactionsController) CreateTransaction(txReq model.CreateTransactionRequest) (*model.Transaction, error) {
	var acct model.Account

	amtToStore, err := json.Marshal(txReq.Money)
	if err != nil {
		return nil, err
	}

	// start transaction
	tx := tc.db.Gorm.Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	args := &model.Transaction{
		AccountID: txReq.AccountID,
		Money:     datatypes.JSON(amtToStore),
		Direction: txReq.Direction,
		Memo:      txReq.Memo,
		State:     model.TransactionStatePending,
		Version:   0,
	}

	// Lock the account row
	err = tx.Set("gorm:query_option", "FOR UPDATE").First(&acct, &model.Account{ID: args.AccountID}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	create := tc.db.Gorm.Create(&args)
	if create.Error != nil {
		return nil, create.Error
	}

	err = executeAccountChanges(tx, args, &acct, false)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return args, nil
}

// ReverseTransaction reverses a transaction
func (tc *TransactionsController) ReverseTransaction(txID string) error {
	transaction := new(model.Transaction)
	acct := new(model.Account)

	tx := tc.db.Gorm.Begin()

	// Lock the tx row
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(transaction, "id = ?", txID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if transaction.State != model.TransactionStateCompleted {
		tx.Rollback()
		return errors.New("cannot reverse transction that did not complete")
	}

	// Lock the account row
	err = tx.Set("gorm:query_option", "FOR UPDATE").First(acct, "id = ?", transaction.AccountID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if transaction.Direction == model.TransactionDirectionDebit {
		transaction.Direction = model.TransactionDirectionCredit
	} else {
		transaction.Direction = model.TransactionDirectionDebit
	}

	err = executeAccountChanges(tx, transaction, acct, true)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func calcAccountBalanceChange(tx *model.Transaction, acct *model.Account) error {
	// Get account balance
	acctBal := new(model.AccountMoney)
	txMoney := new(model.TransactionMoney)

	err := json.Unmarshal(acct.Balance, acctBal)
	if err != nil {
		return err
	}

	err = json.Unmarshal(tx.Money, txMoney)
	if err != nil {
		return err
	}

	if tx.Direction == model.TransactionDirectionDebit {
		if txMoney.Amount > acctBal.Amount {
			return errors.New("insufficient funds")
		}
		acctBal.Amount -= txMoney.Amount
	} else {
		acctBal.Amount += txMoney.Amount
	}

	acctB, err := json.Marshal(acctBal)
	if err != nil {
		return err
	}

	acct.Balance = datatypes.JSON(acctB)

	return nil
}

func updateAccountBalance(gtx *gorm.DB, acct *model.Account) error {
	balanceJSON, err := json.Marshal(acct.Balance)
	if err != nil {
		return err
	}

	// Update account balance
	err = gtx.Model(acct).Update("balance", datatypes.JSON(balanceJSON)).Error
	if err != nil {
		return err
	}

	// Commit the transaction
	err = gtx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func isValidCurrency(txMoney datatypes.JSON, acctMoney datatypes.JSON) bool {
	// Get account balance
	balance := new(model.AccountMoney)
	transaction := new(model.TransactionMoney)

	err := json.Unmarshal(acctMoney, balance)
	if err != nil {
		return false
	}

	err = json.Unmarshal(txMoney, transaction)
	if err != nil {
		return false
	}

	if transaction.Currency != balance.Currency {
		return false
	}

	return true
}
func executeAccountChanges(gtx *gorm.DB, transaction *model.Transaction, account *model.Account, reversing bool) error {
	if reversing {
		transaction.State = model.TransactionStateReversed
	} else {
		if !common.IsCompletableState(*transaction) {
			gtx.Rollback()
			return errors.New("transaction is not in a completable state")
		}
		transaction.State = model.TransactionStateCompleted
	}

	// Update account balance
	err := calcAccountBalanceChange(transaction, account)
	if err != nil {
		gtx.Rollback()
		return err
	}

	// Save transaction status
	err = gtx.Model(*transaction).Update("state", transaction.State).Error
	if err != nil {
		return err
	}

	err = updateAccountBalance(gtx, account)
	if err != nil {
		gtx.Rollback()
		return err
	}

	return nil
}
