package controller

import (
	"github.com/partyscript/bledger/internal/config"
	"github.com/partyscript/bledger/internal/db"

	"github.com/partyscript/bledger/internal/cache"
	"go.uber.org/zap"
)

// Manager is the struct that the constructor implements
type Manager struct {
	Cfg          *config.GlobalConfig
	Transactions TransactionsController
	Accounts     AccountController
}

// NewControllerManager initializes a Manager
func NewControllerManager(
	logger *zap.SugaredLogger,
	cfg *config.GlobalConfig,
	cache *cache.Manager,
	db *db.Manager,
) Manager {
	transactionController := NewTransactionsController(
		logger,
		cfg,
		cache,
		db,
	)

	accountController := NewAccountController(
		logger,
		cfg,
		cache,
		db,
	)

	return Manager{
		Cfg:          cfg,
		Transactions: transactionController,
		Accounts:     accountController,
	}
}
