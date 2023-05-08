package model

import (
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TransactionDirection is an enum for transaction direction
type TransactionDirection string

const (
	// TransactionDirectionDebit is a debit transaction
	TransactionDirectionDebit TransactionDirection = "DEBIT"
	// TransactionDirectionCredit is a credit transaction
	TransactionDirectionCredit TransactionDirection = "CREDIT"
)

// TransactionState is an enum for transaction state
type TransactionState string

const (
	// TransactionStatePending is a pending transaction
	TransactionStatePending TransactionState = "PENDING"
	// TransactionStateCompleted is a completed transaction
	TransactionStateCompleted TransactionState = "COMPLETED"
	// TransactionStateFailed is a failed transaction
	TransactionStateFailed TransactionState = "FAILED"
	// TransactionStateReversed is a reversed transaction
	TransactionStateReversed TransactionState = "REVERSED"
)

// Transaction is the model for a transaction
type Transaction struct {
	gorm.Model  `json:"-"`
	ID          string               `gorm:"primaryKey;uniqueIndex" json:"id"`
	Money       datatypes.JSON       `json:"money"`
	Memo        string               `json:"memo,omitempty"`
	Direction   TransactionDirection `json:"direction"`
	AccountID   string               `json:"account_id"`
	Version     uint                 `json:"version"`
	State       TransactionState     `json:"state" gorm:"default:'PENDING'"`
	ErrorReason string               `json:"error_reason,omitempty"`
}

// TransactionMoney is the model for a transaction money
type TransactionMoney struct {
	Amount   uint64 `binding:"required" gorm:"-" json:"amount"`
	Currency string `binding:"required" gorm:"-" json:"currency"`
}

// CreateTransactionRequest is the model for a transaction create request
type CreateTransactionRequest struct {
	Money     TransactionMoney     `binding:"required" json:"money"`
	Memo      string               `binding:"required" json:"memo,omitempty"`
	Direction TransactionDirection `binding:"required" json:"direction"`
	AccountID string               `binding:"required" json:"account_id"`
}

// BeforeCreate is a method hook that generates a custom sorted id
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = ksuid.New().String()
	return nil
}
