package model

import (
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Account is the model for an account
type Account struct {
	gorm.Model  `json:"-"`
	ID          string         `gorm:"primaryKey;uniqueIndex" json:"id"`
	Balance     datatypes.JSON `json:"balance"`
	Name        string         `binding:"required" json:"name"`
	Description string         `binding:"required" json:"description"`
}

// AccountMoney is the model for an account balance in a currnecy
type AccountMoney struct {
	Amount   uint64 `binding:"required" gorm:"-" json:"amount"`
	Currency string `binding:"required" gorm:"-" json:"currency"`
}

// CreateAccountRequest is the model for an account create request
type CreateAccountRequest struct {
	Name        string `binding:"required" json:"name"`
	Description string `binding:"required" json:"description"`
}

// BeforeCreate is a method hook that generates a custom sorted id
func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = ksuid.New().String()
	return nil
}

// New is a constructor that returns a new instance of Account
func (a *Account) New(id string) *Account {
	return &Account{
		ID: id,
	}

}
