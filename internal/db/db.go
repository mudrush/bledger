package db

import (
	"github.com/partyscript/bledger/internal/model"
	"gorm.io/gorm"
)

// Manager is the struct that the constructor implements
type Manager struct {
	Gorm *gorm.DB
}

// NewDBManager is a constructor that returns a new instance of Manager
// and runs a sql migration
func NewDBManager(gorm *gorm.DB) (*Manager, error) {
	err := migrate(gorm)
	if err != nil {
		return &Manager{}, err
	}
	return &Manager{
		Gorm: gorm,
	}, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Transaction{},
		&model.Account{},
	)
	if err != nil {
		return err
	}
	return nil
}
